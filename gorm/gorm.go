package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	os.Remove("./go_sql_proxy.sqlite3")
	db, err := gorm.Open("sqlite3", "./go_sql_proxy.sqlite3")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	create := `
	create table gorm_test (id integer not null primary key, name text);
	`

	if err := db.Exec(create).Error; err != nil {
		log.Fatal(err)
	}
	if err := db.Exec(`INSERT INTO gorm_test VALUES (1, "1")`).Error; err != nil {
		log.Fatal(err)
	}

	// 存在しないレコードを削除してもエラーにならないことを確認
	if err := db.Exec(`DELETE FROM gorm_test WHERE id = 2`).Error; err != nil {
		log.Fatal(err)
	}

	fmt.Printf("DELETE OK\n")
}
