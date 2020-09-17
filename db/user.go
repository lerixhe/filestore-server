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

// UserSignin 判断数据库密码是否一致
func UserSignin(username, encpwd string) bool {
	sql := "select * from tbl_user where user_name=? limit 1"
	stmt, err := mydb.DBConn().Prepare(sql)
	if err != nil {
		fmt.Printf("Prepare err:%v\n", err)
		return false
	}
	defer stmt.Close()
	rows, err := stmt.Query(username)
	if err != nil {
		fmt.Printf("Query err:%v\n", err)
		return false
	} else if rows == nil {
		fmt.Printf("username no found: %s\n", username)
		rows.Close()
		return false
	}
	defer rows.Close()
	r := mydb.ParseRows(rows)
	if len(r) > 0 && string(r[0]["user_pwd"].([]byte)) == encpwd {
		return true
	}
	return false
}

// UpdateToken 刷新用户的token到数据库
func UpdateToken(username, token string) bool {
	sql := "replace into tbl_user_token(`user_name`,`user_token`)values(?,?)"
	stmt, err := mydb.DBConn().Prepare(sql)
	if err != nil {
		fmt.Printf("Prepare err:%v\n", err)
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, token)
	if err != nil {
		fmt.Printf("Failed to exec replace into token,err:%v\n", err)
		return false
	}
	return true
}

// TokenExist 查询token在数据库是否存在
func TokenExist(username, token string) bool {
	var tokendb string
	sql := "select `user_token` from tbl_user_token where `user_name` = ? limit 1"
	stmt, err := mydb.DBConn().Prepare(sql)
	if err != nil {
		fmt.Printf("Prepare err:%v\n", err)
		return false
	}
	defer stmt.Close()

	row := stmt.QueryRow(username)
	if row == nil {
		return false
	}
	err = row.Scan(&tokendb)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if tokendb != token {
		return false
	}
	return true
}

type UserInfo struct {
	Username     string
	Email        string
	Phone        string
	SignupAt     string
	LastActiveAt string
	Statys       int
}

// GetUserInfo 数据库中查询用户信息
func GetUserInfo(username string) *UserInfo {
	sql := "select `user_name`,`signup_at` from tbl_user where `user_name` = ? limit 1"
	stmt, err := mydb.DBConn().Prepare(sql)
	if err != nil {
		fmt.Printf("Prepare err:%v\n", err)
		return nil
	}
	defer stmt.Close()

	row := stmt.QueryRow(username)
	if row == nil {
		return nil
	}
	user := UserInfo{}
	err = row.Scan(&user.Username, &user.SignupAt)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return &user
}
