package handler

import (
	"encoding/json"
	dblayer "filstore-server/db"
	"filstore-server/meta"
	"filstore-server/util"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

// UploadHandler 区别GET和POST，分别展示页面和接收文件
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w, "internal server error")
			return
		}
		w.Write(data)
	} else if r.Method == http.MethodPost {
		// 接收文件流
		file, head, err := r.FormFile("file")
		if err != nil {
			fmt.Printf("failed to get data,err:%v\n", err)
			return
		}
		defer file.Close()
		fileMeta := meta.FileMeta{
			FileName: head.Filename,
			Location: "./tmp/" + head.Filename,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}
		// 创建新文件句柄
		newFile, err := os.Create(fileMeta.Location)
		if err != nil {
			fmt.Printf("failed to create new file %s,err:%v\n", head.Filename, err)
			return
		}
		defer newFile.Close()
		// 文件流写入新文件
		fileMeta.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("failed to save data into file,err:%v\n", err)
			return
		}
		// 重置文件游标
		newFile.Seek(0, 0)
		fileMeta.FileSha1 = util.FileSha1(newFile)
		meta.UpdateFileMetaDB(fileMeta)

		// TODO:更新用户文件表记录
		r.ParseForm()
		username := r.Form.Get("username")
		if !dblayer.OnUserFileUploadFinished(username, fileMeta.FileSha1, fileMeta.FileName, fileMeta.FileSize) {
			io.WriteString(w, "Upload Failed")
			return
		}

		fmt.Printf("upload file succeed!,total %.2f KB\n", float64(fileMeta.FileSize)/1024)
		// 跳转结果页
		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
	}

}

// UploadSucHandler 上传完成
func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload finished!")
}

// GetFileMetaHandler 获取文件元信息
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	// 将参数Parse到Form中
	r.ParseForm()
	fSha1 := r.Form["filehash"][0]
	fMeta, err := meta.GetFileMetaDB(fSha1)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	data, err := json.Marshal(fMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// DownloadHandler 文件下载
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fSha1 := r.Form.Get("filehash")
	fm := meta.GetFileMeta(fSha1)

	f, err := os.Open(fm.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// 修改http头
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment;filename=\""+fm.FileName+"\"")
	w.Write(data)
}

// FileMetaUpdateHandler 更新元信息
func FileMetaUpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	opType := r.Form.Get("op")
	fileSha1 := r.Form.Get("filehash")
	newFileName := r.Form.Get("filename")

	if opType != "0" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	m := meta.GetFileMeta(fileSha1)
	m.FileName = newFileName
	meta.UpdateFileMeta(m)

	data, err := json.Marshal(m)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// FileDeleteHandler 删除文件
func FileDeleteHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fileSha1 := r.Form.Get("filehash")
	m := meta.GetFileMeta(fileSha1)
	meta.RemoveFileMeta(fileSha1)
	os.Remove(m.FileName)
	w.WriteHeader(http.StatusOK)
}

// FileQueryHandler 处理查询用户文件列表的请求
func FileQueryHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	limitCnt, _ := strconv.Atoi(r.Form.Get("limit"))
	userFiles, err := dblayer.QueryUserFileMetas(username, limitCnt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(userFiles)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// fmt.Println(userFiles)
	w.Write(data)
}
