package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"strings"
	"time"
)

var (
	host     = flag.String("h", "127.0.0.1", "host")
	port     = flag.Int("P", 3306, "port")
	user     = flag.String("u", "test", "user")
	password = flag.String("p", "test", "password")
)

func init() {
	flag.Parse()
}

func main() {
	db, e := sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%d)/information_schema",
			*user, *password, *host, *port))
	if e != nil {
		log.Fatal(e)
		return
	}
	if e := db.Ping(); e != nil {
		log.Fatal(e)
		return
	}
	defer db.Close()
	type status struct {
		questions    uint64
		comSelect    uint64
		comInsert    uint64
		comUpdate    uint64
		comDelete    uint64
		rowsInsert   uint64
		rowsSelect   uint64
		rowsUpdate   uint64
		rowsDelete   uint64
		threadCreate uint64
	}

	handler := func(sta *status, db *sql.DB) {
		ssql := "show global status where variable_name in " +
			"('Questions', 'Com_select', 'Com_update', " +
			"'Com_insert', 'Com_delete', 'Threads_connected', " +
			"'Threads_created', 'Threads_running', " +
			"'Innodb_rows_inserted', 'Innodb_rows_read', " +
			"'Innodb_rows_updated', 'Innodb_rows_deleted')"
		rows, e := db.Query(ssql)
		if e != nil {
			log.Fatal(e)
			os.Exit(1)
		}
		var qps, sel, ins, upd, del, rre, rin, rup, rdel, con, cre, run uint64
		for rows.Next() {
			var (
				key   string
				value uint64
			)
			if err := rows.Scan(&key, &value); err != nil {
				return
			}
			switch key {
			case "Com_delete":
				del = value - sta.comDelete
				sta.comDelete = value
			case "Com_insert":
				ins = value - sta.comInsert
				sta.comInsert = value
			case "Questions":
				qps = value - sta.questions
				sta.questions = value
			case "Com_update":
				upd = value - sta.comUpdate
				sta.comUpdate = value
			case "Com_select":
				sel = value - sta.comSelect
				sta.comSelect = value
			case "Innodb_rows_inserted":
				rin = value - sta.rowsInsert
				sta.rowsInsert = value
			case "Innodb_rows_read":
				rre = value - sta.rowsSelect
				sta.rowsSelect = value
			case "Innodb_rows_updated":
				rup = value - sta.rowsUpdate
				sta.rowsUpdate = value
			case "Innodb_rows_deleted":
				rdel = value - sta.rowsDelete
				sta.rowsDelete = value
			case "Threads_connected":
				con = value
			case "Threads_created":
				cre = value - sta.threadCreate
				sta.threadCreate = value
			case "Threads_running":
				run = value
			default:
			}
		}
		if qps > 100000 {
			return
		}
		timeNow := strings.Split(strings.Fields(time.Now().String())[1], ".")[0]
		fmt.Println(fmt.Sprintf("%s|%5d%6d%6d%7d%7d|%5d%6d%6d%7d|%4d%5d%5d",
			timeNow, ins, upd, del, sel, qps, rin, rup, rdel, rre, run, con, cre))
	}
	state := new(status)
	var t int = 0
	for {
		switch t % 15 {
		case 0:
			fmt.Println(strings.Repeat("_", 80))
			fmt.Println(fmt.Sprintf("--------|%19s%12s| %s |%12s--",
				"--QPS--", " ", "--Innodb Rows Status--", "--Thead"))
			fmt.Println(fmt.Sprintf("  time  |%5s%6s%6s%7s%7s|%5s%6s%6s%7s|%4s%5s%5s",
				"ins", "upd", "del", "sel", "qps", "ins", "upd", "del", "read", "run", "con", "cre"))
			go handler(state, db)
			time.Sleep(time.Second)
			t++
		default:
			go handler(state, db)
			time.Sleep(time.Second)
			t++
		}
	}
}
