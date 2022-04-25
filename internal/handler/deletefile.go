package handler

import (
	"cess-gateway/configs"
	"cess-gateway/internal/chain"
	"cess-gateway/internal/db"
	. "cess-gateway/internal/logger"
	"cess-gateway/tools"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DeletefileHandler(c *gin.Context) {
	var resp = RespMsg{
		Code: http.StatusBadRequest,
		Msg:  "",
	}

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		Err.Sugar().Errorf("%v,%v", c.ClientIP(), err)
		resp.Msg = "bad request"
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	var reqmsg ReqDeleteFileMsg
	err = json.Unmarshal(body, &reqmsg)
	if err != nil {
		Err.Sugar().Errorf("%v,%v", c.ClientIP(), err)
		resp.Msg = "body format error"
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	// Determine if the user has uploaded the file
	key, err := tools.CalcMD5(reqmsg.Filename)
	if err != nil {
		resp.Msg = "invalid filename"
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	resp.Code = http.StatusInternalServerError
	db, err := db.GetDB()
	if err != nil {
		resp.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	fid, err := db.Get(key)
	if err != nil {
		if err.Error() == "leveldb: not found" {
			resp.Code = http.StatusNotFound
			resp.Msg = "This file has not been uploaded"
			c.JSON(http.StatusNotFound, resp)
			return
		} else {
			resp.Msg = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
	}

	//Delete files in cess storage service
	err = chain.DeleteFileOnChain(configs.Confile.AccountSeed, configs.Confile.AccountAddr, fmt.Sprintf("%v", tools.BytesToInt64(fid)))
	if err != nil {
		resp.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	db.Delete(key)
	resp.Code = http.StatusOK
	resp.Msg = "success"
	c.JSON(http.StatusOK, resp)
	return
}
