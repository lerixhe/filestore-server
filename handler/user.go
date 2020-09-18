package handler

import (
	dblayer "filstore-server/db"
	"filstore-server/util"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const (
	pwdSalt   string = "~@3ecu98"
	tokenSalt string = "_tokensalt"
)

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

	encPasswd := util.Sha1([]byte(passwd + pwdSalt))
	if dblayer.UserSignUp(username, encPasswd) {
		w.Write([]byte("SUCCESS"))
	} else {
		w.Write([]byte("FAILED"))
	}
}

// SignInHandler 处理用户登录请求
func SignInHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	passwd := r.Form.Get("password")
	encPasswd := util.Sha1([]byte(passwd + pwdSalt))
	// 校验用户
	if !dblayer.UserSignin(username, encPasswd) {
		w.Write([]byte("FAILED"))
		return
	}
	// 生成token
	token := GenToken(username)
	if !dblayer.UpdateToken(username, token) {
		w.Write([]byte("FAILED"))
		return
	}
	// 由客户端重定向到首页
	// w.Write([]byte("{\"Location\":\"http://" + r.Host + "/static/view/home.html\"}"))
	resp := util.NewResMsg(
		0,
		"OK",
		struct {
			Location string
			Username string
			Token    string
		}{
			"http://" + r.Host + "/static/view/home.html",
			username,
			token,
		},
	)
	w.Write(resp.JSONBytes())
}

// GenToken 生成token
func GenToken(username string) string {
	// 规则：40位字符=md5(username+timestamp+token_salt)+timestamp[:8]
	ts := fmt.Sprintf("%d", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + ts + tokenSalt))
	return tokenPrefix + ts[:8]
}

// UserInfoHandler 处理用户信息请求
func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	user := dblayer.GetUserInfo(username)
	if user == nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	resp := util.NewResMsg(0, "OK", *user)

	w.Write(resp.JSONBytes())
}

// IsTokenValid 验证token有效性
func IsTokenValid(username, token string) bool {
	// 1.判断是否过期
	// 2.判断是否能从数据库查到
	ts := fmt.Sprintf("%d", time.Now().Unix())
	now, err := strconv.Atoi(ts[len(ts)-8:])
	if err != nil {
		fmt.Println(err)
		return false
	}
	then, err := strconv.Atoi(token[32:])
	if err != nil {
		fmt.Println(err)
		return false
	}
	if now-then >= 24*60*60 {
		return false
	}

	if !dblayer.TokenExist(username, token) {
		return false
	}
	return true
}
