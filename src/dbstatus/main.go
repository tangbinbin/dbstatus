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
	host     = flag.String("h", "127.0.0.1:3306", "hosts,多个地址之间,分割")
	user     = flag.String("u", "test", "user")
	password = flag.String("p", "test", "password")
)

func init() {
	flag.Parse()
}

type Info struct {
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

type Server struct {
	conn *sql.DB
	info *Info
}

func main() {
	addrs := strings.Split(*host, ",")
	servers := make(map[string]*Server)
	for _, addr := range addrs {
		db, e := sql.Open("mysql",
			fmt.Sprintf("%s:%s@tcp(%s)/information_schema",
				*user, *password, addr))
		if e != nil {
			log.Fatal(e)
			return
		}
		if e := db.Ping(); e != nil {
			log.Fatal(e)
			return
		}
		server := &Server{
			conn: db,
			info: new(Info),
		}

		servers[addr] = server
	}

	handler := func(addr string, sta *Info, db *sql.DB) {
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
		fmt.Println(fmt.Sprintf("%21s %s|%5d%6d%6d%7d%7d|%5d%6d%6d%7d|%4d%5d%5d",
			addr, timeNow, ins, upd, del, sel, qps, rin, rup, rdel, rre, run, con, cre))
	}
	var (
		t int = 0
		j int = 5
	)
	if len(servers) == 1 {
		j = 15
	}
	for range time.NewTicker(time.Second).C {
		switch t % j {
		case 0:
			fmt.Println("    ", strings.Repeat("_", 98))
			fmt.Println(fmt.Sprintf("                              |%19s%12s| %s |%12s--",
				"--QPS--", " ", "--Innodb Rows Status--", "--Thead"))
			fmt.Println(fmt.Sprintf("          addr          time  |%5s%6s%6s%7s%7s|%5s%6s%6s%7s|%4s%5s%5s",
				"ins", "upd", "del", "sel", "qps", "ins", "upd", "del", "read", "run", "con", "cre"))
			for addr, server := range servers {
				handler(addr, server.info, server.conn)
			}
			if len(servers) > 1 {
				fmt.Println()
			}
			t++
		default:
			for addr, server := range servers {
				handler(addr, server.info, server.conn)
			}
			if len(servers) > 1 {
				fmt.Println()
			}
			t++
		}
	}
}
