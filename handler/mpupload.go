package handler

import (
	"bufio"
	rPool "filstore-server/cache/redis"
	"filstore-server/util"
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"
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
func UploadPartHandler(w http.ResponseWriter, r http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	uploadID := r.Form.Get("uploadid")
	chunkIndex := r.Form.Get("index")

	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	// 获取文件句柄，用于存储分块内容
	fd, err := os.Create("./data/" + uploadID + "/" + chunkIndex)
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
