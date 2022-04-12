package handler

import (
	"cess-httpservice/configs"
	"cess-httpservice/internal/db"
	. "cess-httpservice/internal/logger"
	"cess-httpservice/tools"
	"fmt"
	"strconv"
	"strings"
	"time"

	"encoding/json"

	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/singleflight"
)

// Handler at user registration
func GenerateRandomkeyHandler(c *gin.Context) {
	var resp = RespRandomMsg{
		Code:    http.StatusBadRequest,
		Msg:     "",
		Random1: 0,
		Random2: 0,
	}
	var sf singleflight.Group
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		Err.Sugar().Errorf("%v,%v", c.ClientIP(), err)
		resp.Msg = "bad request"
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	var reqmsg ReqRandomkeyMsg
	err = json.Unmarshal(body, &reqmsg)
	if err != nil {
		Err.Sugar().Errorf("%v,%v", c.ClientIP(), err)
		resp.Msg = "body format error"
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	resp.Code = http.StatusInternalServerError
	randomkey1, randomkey2 := 0, 0
	value := ""
	v, err, _ := sf.Do(reqmsg.Walletaddr, func() (interface{}, error) {
		db, err := db.GetDB()
		if err != nil {
			Err.Sugar().Errorf("%v,%v", c.ClientIP(), err)
			return nil, err
		}
		bytes, err := db.Get([]byte(reqmsg.Walletaddr + "_random"))
		if err != nil && err.Error() != "leveldb: not found" {
			Err.Sugar().Errorf("%v,%v", c.ClientIP(), err)
			return nil, err
		}
		if len(bytes) == 0 {
			randomkey1, randomkey2, value = generateRandoms()
			err = db.Put([]byte(reqmsg.Walletaddr+"_random"), []byte(value))
			if err != nil {
				Err.Sugar().Errorf("%v,%v", c.ClientIP(), err)
				resp.Msg = err.Error()
				c.JSON(http.StatusInternalServerError, resp)
				return nil, err
			}
			return value, err
		} else {
			if len(strings.Split(string(bytes), "#")) != 3 {
				randomkey1, randomkey2, value = generateRandoms()
				err = db.Put([]byte(reqmsg.Walletaddr+"_random"), []byte(value))
				if err != nil {
					Err.Sugar().Errorf("%v,%v", c.ClientIP(), err)
					resp.Msg = err.Error()
					c.JSON(http.StatusInternalServerError, resp)
					return nil, err
				}
				return value, nil
			} else {
				fmt.Println(string(bytes))
				return string(bytes), nil
			}
		}
	})
	if err != nil {
		resp.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	values := strings.Split(v.(string), "#")
	expire, _ := strconv.ParseInt(values[2], 10, 64)
	if time.Since(time.Unix(expire, 0)).Minutes() > configs.RandomValidTime {
		randomkey1, randomkey2, value = generateRandoms()
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
		resp.Code = http.StatusOK
		resp.Msg = "success"
		resp.Random1 = randomkey1
		resp.Random2 = randomkey2
		c.JSON(http.StatusOK, resp)
		return
	}
	resp.Code = http.StatusOK
	resp.Msg = "success"
	resp.Random1, _ = strconv.Atoi(values[0])
	resp.Random2, _ = strconv.Atoi(values[1])
	c.JSON(http.StatusOK, resp)
	return
}

//
func generateRandoms() (int, int, string) {
	for {
		randomkey1 := tools.RandomInRange(100000, 999999)
		randomkey2 := tools.RandomInRange(100000, 999999)
		if randomkey2 != randomkey1 {
			value := fmt.Sprintf("%v", randomkey1) + "#" + fmt.Sprintf("%v", randomkey2) + "#" + fmt.Sprintf("%v", time.Now().Unix())
			return randomkey1, randomkey2, value
		}
	}
}
