package handler

import (
	"cess-httpservice/configs"
	. "cess-httpservice/internal/logger"
	"cess-httpservice/internal/token"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func UpfileHandler(c *gin.Context) {
	var resp = RespMsg{
		Code: http.StatusBadRequest,
		Msg:  "",
		Data: nil,
	}
	usertoken_en := c.PostForm("token")
	bytes, err := token.DecryptToken(usertoken_en)
	if err != nil {
		resp.Msg = "illegal token"
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	var usertoken token.TokenMsgType
	err = json.Unmarshal(bytes, &usertoken)
	if err != nil {
		resp.Msg = "token format error"
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	fmt.Println("token:", usertoken_en)
	content_length := c.Request.ContentLength
	if content_length <= 0 {
		resp.Msg = "empty file"
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	file_p, err := c.FormFile("file")
	if err != nil {
		resp.Msg = "not upload file request"
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	file_c, _, _ := c.Request.FormFile("file")
	userpath := filepath.Join(configs.FileCacheDir, fmt.Sprintf("%v", usertoken.Userid))
	_, err = os.Stat(userpath)
	if err != nil {
		err = os.MkdirAll(userpath, os.ModeDir)
		if err != nil {
			resp.Code = http.StatusInternalServerError
			resp.Msg = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
	}

	fpath := filepath.Join(userpath, file_p.Filename)
	_, err = os.Stat(fpath)
	if err == nil {
		resp.Msg = "duplicate filename"
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	resp.Code = http.StatusInternalServerError
	f, err := os.Create(fpath)
	if err != nil {
		Err.Sugar().Errorf("create file fail:%v\n", err)
		resp.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	defer f.Close()
	buf := make([]byte, 2*1024*1024)
	for {
		n, err := file_c.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			resp.Code = http.StatusGatewayTimeout
			resp.Msg = "upload failed due to network issues"
			c.JSON(http.StatusGatewayTimeout, resp)
			return
		}
		if n == 0 {
			continue
		}
		f.Write(buf[:n])
	}
	resp.Code = http.StatusOK
	resp.Msg = "success"
	c.JSON(http.StatusOK, resp)

	//TODO: Upload files to cess storage system

	return
}
