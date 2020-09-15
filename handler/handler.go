package handler

import (
	"io"
	"io/ioutil"
	"net/http"
)

// UploadHandler 区别GET和POST，分别展示页面和接收文件
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w, "internal server error")
			return
		}
		io.WriteString(w, string(data))

	} else if r.Method == http.MethodPost {

	}

}
