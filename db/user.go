package db

import (
	mydb "filstore-server/db/mysql"
	"fmt"
)

// UserSignUp 通过用户名和密码完成user表的注册操作
func UserSignUp(userName, passwd string) bool {
	sql := "insert ignore into tbl_user(`user_name`,`user_pwd`)values(?,?)"
	stmt, err := mydb.DBConn().Prepare(sql)
	if err != nil {
		fmt.Printf("Failed to insert user,err:%v\n", err)
		return false
	}
	defer stmt.Close()
	re, err := stmt.Exec(userName, passwd)
	if err != nil {
		fmt.Printf("Failed to exec insert user,err:%v\n", err)
		return false
	}
	rf, err := re.RowsAffected()
	if err != nil {
		fmt.Printf("Failed re.RowsAffected,err:%v\n", err)
		return false
	} else if rf == 0 {
		fmt.Printf("user %s had been signup before!\n", userName)
		return false
	}
	return true
}
