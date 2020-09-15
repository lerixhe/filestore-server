package handler

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
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
		// 创建新文件句柄
		newFile, err := os.Create("./tmp/" + head.Filename)
		if err != nil {
			fmt.Printf("failed to create new file %s,err:%v\n", head.Filename, err)
			return
		}
		defer newFile.Close()
		// 文件流写入新文件
		len, err := io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("failed to save data into file,err:%v\n", err)
			return
		}
		fmt.Printf("upload file succeed!,total %.2f KB\n", float64(len)/1024)
		// 跳转结果页
		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
	}

}

// UploadSucHandler 上传完成
func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload finished!")
}
