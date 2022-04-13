package handler

import (
	"cess-httpservice/internal/db"
	"cess-httpservice/internal/token"
	"cess-httpservice/tools"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func DownfileHandler(c *gin.Context) {
	var resp = RespMsg{
		Code: http.StatusBadRequest,
		Msg:  "",
	}
	usertoken := c.Query("token")
	filename := c.Query("filename")
	if usertoken == "" || filename == "" {
		resp.Msg = "token or filename is empty"
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	//Parse token
	bytes, err := token.DecryptToken(usertoken)
	if err != nil {
		resp.Msg = "invalid token"
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	var token_de token.TokenMsgType
	err = json.Unmarshal(bytes, &token_de)
	if err != nil {
		resp.Msg = "token format error"
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	if time.Now().Unix() > token_de.Expire {
		resp.Code = http.StatusForbidden
		resp.Msg = "token expired"
		c.JSON(http.StatusForbidden, resp)
		return
	}

	//Determine if the user has uploaded the file
	key, err := tools.CalcMD5(fmt.Sprintf("%v", token_de.Userid) + filename)
	if err != nil {
		resp.Msg = "invalid filename"
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	db, err := db.GetDB()
	if err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	ok, err := db.Has(key)
	if err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	if !ok {
		resp.Code = http.StatusNotFound
		resp.Msg = "This file has not been uploaded"
		c.JSON(http.StatusNotFound, resp)
		return
	}

	//TODO: Download the file from the scheduler service

	//c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filehash))
	c.Writer.Header().Add("inline", fmt.Sprintf("inline; filename=%v", filename))
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.File("/usr/local/cess/cache/test/goo.mod")
}
