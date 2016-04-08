package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"runtime"
	"sort"
	"strings"
	"time"
)

var (
	host         = flag.String("h", "127.0.0.1:3306", "hosts,多个地址之间,分割")
	user         = flag.String("u", "test", "user")
	password     = flag.String("p", "test", "password")
	length   int = 0
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
	byteReceive  uint64
	byteSent     uint64
}

type State struct {
	qps  uint64
	sel  uint64
	ins  uint64
	upd  uint64
	del  uint64
	rin  uint64
	rre  uint64
	rup  uint64
	rdel uint64
	con  uint64
	cre  uint64
	run  uint64
	recv uint64
	send uint64
}

type Server struct {
	id      int
	addr    string
	conn    *sql.DB
	info    *Info
	state   *State
	timeNow string
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	addrs := strings.Split(*host, ",")
	servers := make(map[int]*Server)
	id := 0
	for _, addr := range addrs {
		db, e := sql.Open("mysql",
			fmt.Sprintf(
				"%s:%s@tcp(%s)/information_schema",
				*user, *password, addr,
			),
		)
		if e != nil {
			log.Fatal(e)
			return
		}
		if e := db.Ping(); e != nil {
			log.Fatal(e)
			return
		}
		server := &Server{
			id:    id,
			addr:  addr,
			conn:  db,
			info:  &Info{},
			state: &State{},
		}

		servers[id] = server
		id++

		if len(addr) > length {
			length = len(addr)
		}
	}
	sk := make([]int, 0)
	for k, _ := range servers {
		sk = append(sk, k)
	}
	sort.Ints(sk)

	go func() {
		for range time.NewTicker(time.Second).C {
			for _, k := range sk {
				getInfo(servers[k])
			}
		}
	}()
	time.Sleep(2 * time.Second)
	i, j := 0, 5
	if len(addrs) == 1 {
		j = 15
	}
	for range time.NewTicker(time.Second).C {
		if i%j == 0 {
			fmt.Println("\033[0;43;1m" +
				strings.Repeat(" ", length+9) +
				"|           --QPS--          | --Innodb Rows Status-- |  --Thread--  | --kbytes-- \033[0m")
			fmt.Println("\033[0;32;1m" +
				strings.Repeat(" ", length-4) +
				"addr     time|  ins  upd  del   sel    qps|  ins   upd   del   read| run  con  cre|  recv  send")
		}
		for _, k := range sk {
			echoState(servers[k])
		}
		if j == 5 {
			fmt.Println()
		}
		i++
	}
}

func echoState(s *Server) {
	fmt.Println(
		fmt.Sprintf("\033[0;34;1m%s \033[0;33;1m%s\033[0m|%5d%5d%5d%6d%7d|%5d%6d%6d%7d|%4d%5d%5d|%6d%6d",
			tolen(s.addr), s.timeNow, s.state.ins, s.state.upd, s.state.del, s.state.sel,
			s.state.qps, s.state.rin, s.state.rup, s.state.rdel, s.state.rre,
			s.state.run, s.state.con, s.state.cre, s.state.recv, s.state.send,
		),
	)
}

func tolen(s string) string {
	if len(s) < length {
		return strings.Repeat(" ", length-len(s)) + s
	}
	return s
}

func getInfo(s *Server) {
	ssql := "show global status where variable_name in " +
		"('Questions', 'Com_select', 'Com_update', " +
		"'Com_insert', 'Com_delete', 'Threads_connected', " +
		"'Threads_created', 'Threads_running', " +
		"'Innodb_rows_inserted', 'Innodb_rows_read', " +
		"'Innodb_rows_updated', 'Innodb_rows_deleted', " +
		"'Bytes_received','Bytes_sent')"
	rows, e := s.conn.Query(ssql)
	s.timeNow = strings.Split(strings.Fields(time.Now().String())[1], ".")[0]
	if e != nil {
		log.Fatal(e)
	}
	for rows.Next() {
		var (
			key   string
			value uint64
		)
		if err := rows.Scan(&key, &value); err != nil {
			return
		}
		switch key {
		case "Bytes_received":
			s.state.recv = (value - s.info.byteReceive) / 1000
			s.info.byteReceive = value
		case "Bytes_sent":
			s.state.send = (value - s.info.byteSent) / 1000
			s.info.byteSent = value
		case "Com_delete":
			s.state.del = value - s.info.comDelete
			s.info.comDelete = value
		case "Com_insert":
			s.state.ins = value - s.info.comInsert
			s.info.comInsert = value
		case "Questions":
			s.state.qps = value - s.info.questions
			s.info.questions = value
		case "Com_update":
			s.state.upd = value - s.info.comUpdate
			s.info.comUpdate = value
		case "Com_select":
			s.state.sel = value - s.info.comSelect
			s.info.comSelect = value
		case "Innodb_rows_inserted":
			s.state.rin = value - s.info.rowsInsert
			s.info.rowsInsert = value
		case "Innodb_rows_read":
			s.state.rre = value - s.info.rowsSelect
			s.info.rowsSelect = value
		case "Innodb_rows_updated":
			s.state.rup = value - s.info.rowsUpdate
			s.info.rowsUpdate = value
		case "Innodb_rows_deleted":
			s.state.rdel = value - s.info.rowsDelete
			s.info.rowsDelete = value
		case "Threads_connected":
			s.state.con = value
		case "Threads_created":
			s.state.cre = value - s.info.threadCreate
			s.info.threadCreate = value
		case "Threads_running":
			s.state.run = value
		default:
		}
	}
}
