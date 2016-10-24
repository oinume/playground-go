package main

import (
	"database/sql"
	"database/sql/driver"
	"log"
	"fmt"
	"os"

	"github.com/shogo82148/go-sql-proxy"
	"github.com/mattn/go-sqlite3"
)

var _ = fmt.Print

func main() {
	fmt.Printf("Remove\n")
	os.Remove("./go_sql_proxy.sqlite3")
	db, err := open("./go_sql_proxy.sqlite3", true)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	fmt.Printf("db = %v\n", db)

	create := `
	create table test (id integer not null primary key, name text);
	delete from test;
	`

	result, err := db.Exec(create)
	if err != nil {
		log.Fatal(err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("result = %v\n", affected)
}

func open(dsn string, useProxy bool) (*sql.DB, error) {
	if useProxy {
		hooks := &proxy.Hooks{
			PreExec: func(stmt *proxy.Stmt, args []driver.Value) (interface{}, error) {
				fmt.Printf("PreExec: stmt=%v\n", stmt.QueryString)
				//panic(stmt.QueryString)
				return nil, nil
			},
			PreQuery: func(stmt *proxy.Stmt, args []driver.Value) (interface{}, error) {
				return nil, nil
			},
			PreBegin: func(_ *proxy.Conn) (interface{}, error) {
				return nil, nil
			},
			PreCommit: func(_ *proxy.Tx) (interface{}, error) {
				return nil, nil
			},
			PreRollback: func(_ *proxy.Tx) (interface{}, error) {
				return nil, nil
			},
		}
		sql.Register("sqlite3-trace", proxy.NewProxy(&sqlite3.SQLiteDriver{}, hooks))
		return sql.Open("sqlite3-trace", dsn)
	} else {
		return sql.Open("sqlite3", dsn)
	}
}
