package handler

import (
	"cess-gateway/configs"
	"cess-gateway/internal/chain"
	"cess-gateway/internal/db"
	. "cess-gateway/internal/logger"
	"cess-gateway/internal/token"
	"cess-gateway/tools"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func DeletefileHandler(c *gin.Context) {
	var resp = RespMsg{
		Code: http.StatusUnauthorized,
		Msg:  Status_401_token,
	}
	// token
	htoken := c.Request.Header.Get("Authorization")
	if htoken == "" {
		Err.Sugar().Errorf("[%v] head missing token", c.ClientIP())
		c.JSON(http.StatusUnauthorized, resp)
		return
	}

	bytes, err := token.DecryptToken(htoken)
	if err != nil {
		Err.Sugar().Errorf("[%v] [%v] DecryptToken error", c.ClientIP(), htoken)
		c.JSON(http.StatusUnauthorized, resp)
		return
	}

	var usertoken token.TokenMsgType
	err = json.Unmarshal(bytes, &usertoken)
	if err != nil {
		Err.Sugar().Errorf("[%v] [%v] token format error", c.ClientIP(), htoken)
		c.JSON(http.StatusUnauthorized, resp)
		return
	}

	if time.Now().Unix() >= usertoken.ExpirationTime {
		Err.Sugar().Errorf("[%v] [%v] token expired", c.ClientIP(), usertoken.Mailbox)
		resp.Msg = Status_401_expired
		c.JSON(http.StatusUnauthorized, resp)
		return
	}

	resp.Code = http.StatusBadRequest
	resp.Msg = Status_400_default
	filename := c.Param("filename")
	if filename == "" {
		Err.Sugar().Errorf("[%v] [%v] no file name", c.ClientIP(), htoken)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	key, err := tools.CalcMD5(usertoken.Mailbox + filename)
	if err != nil {
		Err.Sugar().Errorf("[%v] [%v] %v", c.ClientIP(), usertoken.Mailbox, err)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	resp.Code = http.StatusInternalServerError
	resp.Msg = Status_500_db
	db, err := db.GetDB()
	if err != nil {
		Err.Sugar().Errorf("[%v] [%v] %v", c.ClientIP(), usertoken.Mailbox, err)
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
			Err.Sugar().Errorf("[%v] [%v] %v", c.ClientIP(), usertoken.Mailbox, err)
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
	}

	//Delete files in cess storage service
	err = chain.DeleteFileOnChain(configs.Confile.AccountSeed, configs.Confile.AccountAddr, fmt.Sprintf("%v", tools.BytesToInt64(fid)))
	if err != nil {
		Err.Sugar().Errorf("[%v] [%v] %v", c.ClientIP(), usertoken.Mailbox, err)
	}
	db.Delete(key)
	db.Delete(fid)
	resp.Code = http.StatusOK
	resp.Msg = "success"
	c.JSON(http.StatusOK, resp)
	return
}
