package handler

import (
	"cess-httpservice/internal/chain"
	"cess-httpservice/internal/db"
	"cess-httpservice/internal/token"
	"cess-httpservice/tools"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func FilelistHandler(c *gin.Context) {
	var resp = RespMsg{
		Code: http.StatusBadRequest,
		Msg:  "",
	}
	usertoken := c.Query("token")
	if usertoken == "" {
		resp.Msg = "token is empty"
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	// Parse token
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
	resp.Data = fnames
	c.JSON(http.StatusOK, resp)
}
