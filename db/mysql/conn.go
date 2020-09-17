package mysql

import (
	"database/sql"
	"fmt"
	"os"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	var err error
	// 注意一定不要写冒号
	db, err = sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/fileserver?charset=utf8")
	if err != nil {
		fmt.Printf("Failed to open mysql,err:%v", err)
		os.Exit(1)
	}
	db.SetMaxOpenConns(1000)
	err = db.Ping()
	if err != nil {
		fmt.Printf("Failed to connect to mysql,err:%v", err)
		os.Exit(1)
	}
}

// DBConn 返回数据库连接
func DBConn() *sql.DB {
	return db
}

// ParseRows 将row转为map切片
func ParseRows(rows *sql.Rows) []map[string]interface{} {
	// 总体思路：
	// 1.利用scan，扫描rows到二维数组
	// 2.遍历二维数组的数据转为map切片形式

	// 获取列名
	columns, _ := rows.Columns()
	// 构造二维切片
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for j := range values {
		scanArgs[j] = &values[j]
	}
	// 构造记录map
	record := make(map[string]interface{})
	records := make([]map[string]interface{}, 0)
	// 遍历每一行的数据
	for rows.Next() {
		//扫描1行数据到到二维切片中，其中，第一维是列，第二维是行
		err := rows.Scan(scanArgs...)
		if err != nil {
			panic(err)
		}
		// 遍历二维切片中的有数据的行，并将数据取出存入record的map中
		for i, col := range values {
			if col != nil {
				record[columns[i]] = col
			}
		}
		records = append(records, record)
	}
	return records
}
