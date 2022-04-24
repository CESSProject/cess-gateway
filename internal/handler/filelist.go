package handler

import (
	"cess-httpservice/internal/chain"
	"cess-httpservice/internal/db"
	. "cess-httpservice/internal/logger"
	"cess-httpservice/internal/token"
	"cess-httpservice/tools"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func FilelistHandler(c *gin.Context) {
	var resp = RespMsg{
		Code: http.StatusUnauthorized,
		Msg:  Status_401_token,
	}
	// token
	htoken := c.Request.Header.Get("Authorization")
	fmt.Println(htoken)
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

	resp.Code = http.StatusInternalServerError
	data, err := chain.GetFilelistInfo(token_de.Walletaddr)
	if err != nil {
		resp.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	var fnames = make([]string, 0)
	db, err := db.GetDB()
	if err != nil {
		resp.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	for i := 0; i < len(data); i++ {
		s := ""
		for j := 0; j < len(data[i]); j++ {
			temp := fmt.Sprintf("%c", data[i][j])
			s += temp
		}
		fid, err := strconv.ParseInt(s, 10, 64)
		if err == nil {
			b, err := db.Get(tools.Int64ToBytes(fid))
			if err == nil {
				fnames = append(fnames, string(b))
			}
		}
	}
	resp.Code = http.StatusOK
	resp.Msg = "success"
	//resp.Data = fnames
	c.JSON(http.StatusOK, resp)
}
