package db

import (
	mydb "filstore-server/db/mysql"
	"fmt"
	"time"
)

// UserFile 用户文件信息
type UserFile struct {
	UserName    string
	FileHash    string
	FileName    string
	FileSize    int64
	UploadAt    string
	LastUpdated string
	Status      string
}

// OnUserFileUploadFinished 将用户的文件信息存入数据库
func OnUserFileUploadFinished(username, filehash, filename string, filesize int64) bool {
	sql := "insert ignore into tbl_user_file(`user_name`,`file_sha1`,`file_size`,`file_name`,`upload_at`)values(?,?,?,?,?)"
	stmt, err := mydb.DBConn().Prepare(sql)
	if err != nil {
		fmt.Print(err)
		return false
	}
	defer stmt.Close()
	_, err = stmt.Exec(username, filehash, filesize, filename, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		return false
	}
	return true

}

// QueryUserFileMetas 从数据库中查询用户的文件列表
func QueryUserFileMetas(username string, limit int) ([]UserFile, error) {
	sql := "select file_sha1,file_name,file_size,upload_at,last_update from tbl_user_file where user_name = ? limit ?"
	stmt, err := mydb.DBConn().Prepare(sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(username, limit)
	if err != nil {
		return nil, err
	}
	var userfiles []UserFile
	for rows.Next() {
		var userfile UserFile
		err := rows.Scan(&userfile.FileHash, &userfile.FileName, &userfile.FileSize, &userfile.UploadAt, &userfile.LastUpdated)
		if err != nil {
			fmt.Println(err)
			// 使用break，继续处理下一行
			break
		}
		userfiles = append(userfiles, userfile)
	}
	return userfiles, nil
}
