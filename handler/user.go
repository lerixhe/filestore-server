package handler

import (
	dblayer "filstore-server/db"
	"filstore-server/util"
	"io/ioutil"
	"net/http"
)

const pwdSalt string = "~@3ecu98"

// SignupHandler 处理用户注册请求
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data, err := ioutil.ReadFile("./static/view/signup.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
		return
	}
	r.ParseForm()
	username := r.Form.Get("username")
	passwd := r.Form.Get("password")

	// 简单校验
	if len(username) < 3 || len(passwd) < 5 {
		w.Write([]byte("Invalid parameter"))
		return
	}

	enc_passwd := util.Sha1([]byte(passwd + pwdSalt))
	if dblayer.UserSignUp(username, enc_passwd) {
		w.Write([]byte("SUCCESS"))
	} else {
		w.Write([]byte("FAILED"))
	}
}
