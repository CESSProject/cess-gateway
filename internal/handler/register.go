package handler

import (
	. "cess-httpservice/internal/logger"
	"cess-httpservice/internal/token"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Handler at user registration
func GenerateAccessTokenHandler(c *gin.Context) {
	var resp = RespMsg{
		Code: 1,
		Msg:  "",
		Data: nil,
	}
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		Err.Sugar().Errorf("%v,%v", c.ClientIP(), err)
		resp.Msg = "bad request"
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	var reqmsg RegistrationReq
	err = json.Unmarshal(body, &reqmsg)
	if err != nil {
		Err.Sugar().Errorf("%v,%v", c.ClientIP(), err)
		resp.Msg = "body format error"
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	fmt.Println(reqmsg)
	//TODO: Query the block information to determine whether the wallet address is consistent

	//TODO: Generate user token and store to database
	expire := time.Now().Add(time.Hour * 24 * 7).Unix()
	tk, err := token.GetToken(reqmsg.Walletaddr, reqmsg.Blocknumber, expire)
	if err != nil {
		Err.Sugar().Errorf("%v,%v", c.ClientIP(), err)
		resp.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	resp.Code = 200
	resp.Msg = "success"
	resp.Data = tk
	c.JSON(http.StatusOK, resp)
	return
}
