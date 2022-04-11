package handler

import (
	"cess-httpservice/configs"
	"cess-httpservice/internal/db"
	. "cess-httpservice/internal/logger"
	"cess-httpservice/internal/token"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Handler at user registration
func GenerateAccessTokenHandler(c *gin.Context) {
	var resp = RespMsg{
		Code: http.StatusBadRequest,
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
	var reqmsg ReqRegistrationMsg
	err = json.Unmarshal(body, &reqmsg)
	if err != nil {
		Err.Sugar().Errorf("%v,%v", c.ClientIP(), err)
		resp.Msg = "body format error"
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	//TODO: Query block information

	resp.Code = http.StatusInternalServerError
	db, err := db.GetDB()
	if err != nil {
		Err.Sugar().Errorf("%v,%v", c.ClientIP(), err)
		resp.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	bytes, err := db.Get([]byte(reqmsg.Walletaddr + "_random"))
	if err != nil {
		Err.Sugar().Errorf("%v,%v", c.ClientIP(), err)
		resp.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	value := strings.Split(string(bytes), "#")
	if len(value) != 3 {
		db.Delete([]byte(reqmsg.Walletaddr + "_random"))
		Err.Sugar().Errorf("%v,%v", c.ClientIP(), err)
		resp.Msg = "Please get the random number again (valid within 5 minutes)"
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	randomExpire, err := strconv.Atoi(value[2])
	if time.Since(time.Unix(int64(randomExpire), 0)).Minutes() > configs.RandomValidTime {
		db.Delete([]byte(reqmsg.Walletaddr + "_random"))
		Err.Sugar().Errorf("%v,%v", c.ClientIP(), err)
		resp.Code = http.StatusForbidden
		resp.Msg = "Please get the random number again (valid within 5 minutes)"
		c.JSON(http.StatusForbidden, resp)
		return
	}

	random2Local, err := strconv.Atoi(value[1])
	if time.Since(time.Unix(int64(randomExpire), 0)).Minutes() > configs.RandomValidTime {
		db.Delete([]byte(reqmsg.Walletaddr + "_random"))
		Err.Sugar().Errorf("%v,%v", c.ClientIP(), err)
		resp.Msg = "Please get the random number again (valid within 5 minutes)"
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	//TODO: Judgment random number1
	if reqmsg.Random2 != random2Local {
		Err.Sugar().Errorf("%v,%v", c.ClientIP(), err)
		resp.Code = http.StatusForbidden
		resp.Msg = "Authentication failed"
		c.JSON(http.StatusForbidden, resp)
		return
	}

	//Generate user token
	expire := time.Now().Add(time.Hour * 24 * 7).Unix()
	tk, err := token.GetToken(reqmsg.Walletaddr, reqmsg.Blocknumber, expire)
	if err != nil {
		Err.Sugar().Errorf("%v,%v", c.ClientIP(), err)
		resp.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	//store token to database
	err = db.Put([]byte(reqmsg.Walletaddr+"_token"), []byte(tk))
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
