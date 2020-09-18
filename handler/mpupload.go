package handler

import (
	"bufio"
	rPool "filstore-server/cache/redis"
	dblayer "filstore-server/db"
	"filstore-server/util"
	"fmt"
	"math"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
)

// MultipartUploadInfo 分块信息
type MultipartUploadInfo struct {
	FileHash   string
	FileSize   int
	UploadID   string
	ChunkSize  int
	ChunkCount int
}

// InitialMultipartUploadHandler 初始化分块上传
func InitialMultipartUploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filesize, err := strconv.Atoi(r.Form.Get("filesize"))
	if err != nil {
		w.Write(util.GenSimpleResStream(-1, "params invalid"))
		return
	}

	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	upInfo := MultipartUploadInfo{
		FileHash:   filehash,
		FileSize:   filesize,
		UploadID:   username + fmt.Sprintf("%x", time.Now().UnixNano()),
		ChunkSize:  5 * 1024 * 1024, //5MB
		ChunkCount: int(math.Ceil(float64(filesize) / 5 * 1024 * 1024)),
	}
	rConn.Do("HSET", "MP_"+upInfo.UploadID, "chunkcount", upInfo.ChunkCount)
	rConn.Do("HSET", "MP_"+upInfo.UploadID, "filehash", upInfo.FileHash)
	rConn.Do("HSET", "MP_"+upInfo.UploadID, "filesize", upInfo.FileSize)

	w.Write(util.GenSimpleResStream(0, "OK"))
}

// UploadPartHandler 上传文件分块
func UploadPartHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	// username := r.Form.Get("username")
	uploadID := r.Form.Get("uploadid")
	chunkIndex := r.Form.Get("index")

	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	// 获取文件句柄，用于存储分块内容
	fpath := "./data/" + uploadID + "/" + chunkIndex
	os.MkdirAll(path.Dir(fpath), 0744)
	fd, err := os.Create(fpath)
	if err != nil {
		w.Write(util.GenSimpleResStream(-1, "Upload part failed:"+err.Error()))
		return
	}
	defer fd.Close()

	bw := bufio.NewWriterSize(fd, 1024*1024)
	_, err = bw.ReadFrom(r.Body)
	if err != nil {
		w.Write(util.GenSimpleResStream(-1, "Upload part failed:"+err.Error()))
		return
	}

	rConn.Do("HSET", "MP_"+uploadID, "chkidx_"+chunkIndex, 1)
	w.Write(util.GenSimpleResStream(0, "OK"))

}

// CompleteUploadHandler 通知上传合并
func CompleteUploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	uploadID := r.Form.Get("uploadid")
	filehash := r.Form.Get("filehash")
	filesize := r.Form.Get("filesize")
	filename := r.Form.Get("filename")

	rConn := rPool.RedisPool().Get()
	defer rConn.Close()
	data, err := redis.Values(rConn.Do("HGETALL", "MP_"+uploadID))
	if err != nil {
		w.Write(util.GenSimpleResStream(-1, "complete upload failed:"+err.Error()))
		return
	}
	totalCount, ChunkCount := 0, 0
	for i := 0; i < len(data); i += 2 {
		k := string(data[i].([]byte))
		v := string(data[i+1].([]byte))
		if k == "chunkcount" {
			totalCount, _ = strconv.Atoi(v)
		} else if strings.HasPrefix(k, "chkidx_") && v == "1" {
			ChunkCount++
		}
	}
	if totalCount != ChunkCount {
		w.Write(util.GenSimpleResStream(-2, "invalid request"))
		return
	}
	// 更新唯一文件表和用户文件表
	fsize, _ := strconv.Atoi(filesize)
	dblayer.OnFileUploadFinished(filehash, filename, "", int64(fsize))
	dblayer.OnUserFileUploadFinished(username, filehash, filename, int64(fsize))
	w.Write(util.GenSimpleResStream(0, "OK"))

}
