// the package operate  to mysql database
package conDB

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func Conn() {
	fmt.Println("begin to mysql===============")
	db, err := sql.Open("mysql", "root:123456@tcp(localhost:3306)/test?charset=utf8")
	checkErr(err)

	//查询数据
	rows, err := db.Query("SELECT * FROM file_table")

	for rows.Next() {
		var ID int
		var USERNAME string
		var FILESTATUS string
		var UPTIME string
		var FILENAME string
		var PLAYNAME string

		err = rows.Scan(&USERNAME, &ID, &FILESTATUS, &UPTIME, &FILENAME, &PLAYNAME)
		checkErr(err)

		fmt.Println(ID)
		fmt.Println(USERNAME)
		fmt.Println(FILESTATUS)
		fmt.Println(UPTIME)
		fmt.Println(FILENAME)
		fmt.Println(PLAYNAME)
		fmt.Println("================================================================")
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
