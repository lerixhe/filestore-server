package db

import (
	"database/sql"
	mydb "filstore-server/db/mysql"
	"fmt"
)

// OnFileUploadFinished 上传完成,保存meta
func OnFileUploadFinished(filehash, filename, fileaddr string, filesize int64) bool {
	stmt, err := mydb.DBConn().Prepare("insert ignore into tbl_file(`file_sha1`,`file_name`,`file_size`,`file_addr`,`status`) values(?,?,?,?,1)")
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

type TableFile struct {
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
}

// GetFileMeta 从DB获取meta
func GetFileMeta(fileHash string) (*TableFile, error) {
	stmt, err := mydb.DBConn().Prepare("select file_sha1,file_name,file_addr,file_size from tbl_file where file_sha1=? and status=1 limit 1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	var f TableFile
	err = stmt.QueryRow(fileHash).Scan(&f.FileHash, &f.FileName, &f.FileAddr, &f.FileSize)
	if err != nil {
		return nil, err
	}
	return &f, nil
}
