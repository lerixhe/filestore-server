package meta

import (
	mydb "filstore-server/db"
)

// FileMeta 文件元信息结构
type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)
}

// UpdateFileMeta :新增或更改文件元信息
func UpdateFileMeta(fmeta FileMeta) {
	fileMetas[fmeta.FileSha1] = fmeta
}

// UpdateFileMetaDB :新增或更改数据库中的文件元信息
func UpdateFileMetaDB(fmeta FileMeta) bool {
	return mydb.OnFileUploadFinished(fmeta.FileSha1, fmeta.FileName, fmeta.Location, fmeta.FileSize)
}

// GetFileMeta :获取文件元信息
func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}

// GetFileMetaDB :从数据库获取文件元信息
func GetFileMetaDB(fileSha1 string) (FileMeta, error) {
	var fmeta FileMeta
	ft, err := mydb.GetFileMeta(fileSha1)
	if err != nil {
		return fmeta, err
	}
	fmeta.FileSha1 = fileSha1
	fmeta.FileName = ft.FileName.String
	fmeta.FileSize = ft.FileSize.Int64
	fmeta.Location = ft.FileAddr.String
	return fmeta, nil
}

// RemoveFileMeta :删除文件元信息,暂不考虑线程安全
func RemoveFileMeta(fileSha1 string) {
	delete(fileMetas, fileSha1)
}
