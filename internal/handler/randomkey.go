package handler

import (
	"cess-httpservice/internal/db"
	. "cess-httpservice/internal/logger"
	"cess-httpservice/tools"
	"fmt"
	"time"

	"encoding/json"

	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler at user registration
func GenerateRandomkeyHandler(c *gin.Context) {
	var resp = RespranomMsg{
		Code:    1,
		Msg:     "",
		Random1: 0,
		Random2: 0,
	}
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		Err.Sugar().Errorf("%v,%v", c.ClientIP(), err)
		resp.Msg = "bad request"
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	var reqmsg GetRandomkeyReq
	err = json.Unmarshal(body, &reqmsg)
	if err != nil {
		Err.Sugar().Errorf("%v,%v", c.ClientIP(), err)
		resp.Msg = "body format error"
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	randomkey1 := 0
	randomkey2 := 0
	for {
		randomkey1 = tools.RandomInRange(100000, 999999)
		randomkey2 = tools.RandomInRange(100000, 999999)
		if randomkey2 != randomkey1 {
			break
		}
	}
	db, err := db.GetDB()
	if err != nil {
		Err.Sugar().Errorf("%v,%v", c.ClientIP(), err)
		resp.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	value := fmt.Sprintf("%v", randomkey1) + "#" + fmt.Sprintf("%v", randomkey2) + "#" + fmt.Sprintf("%v", time.Now().Unix())
	err = db.Put([]byte(reqmsg.Walletaddr+"_random"), []byte(value))
	if err != nil {
		Err.Sugar().Errorf("%v,%v", c.ClientIP(), err)
		resp.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	resp.Code = 200
	resp.Msg = "success"
	resp.Random1 = randomkey1
	resp.Random2 = randomkey2
	c.JSON(http.StatusOK, resp)
	return
}
