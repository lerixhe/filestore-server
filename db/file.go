package db

import (
	mydb "filstore-server/db/mysql"
	"fmt"
)

// OnFileUploadFinished 上传完成,保存meta
func OnFileUploadFinished(filehash, filename, fileaddr string, filesize int64) bool {
	stmt, err := mydb.DBConn().Prepare("insert ignore into tbl_file (`file_sha1`,`file_name`,`file_size`,`file_addr`,`status`) values(?,?,?,?,1)")
	if err != nil {
		fmt.Printf("Failed to prepare statement,err:%v\n", err)
		return false
	}
	defer stmt.Close()
	re, err := stmt.Exec(filehash, filename, filesize, fileaddr)
	if err != nil {
		fmt.Printf("Failed to exec statement,err:%v\n", err)
		return false
	}
	if rf, err := re.RowsAffected(); err == nil {
		if rf == 0 {
			fmt.Printf("Warning:File with hash %s had been uploaded before\n", filehash)
		}
		return true
	}
	fmt.Printf("err:%v\n", err)
	return false
}
