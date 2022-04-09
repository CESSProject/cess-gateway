package handler

import (
	. "cess-httpservice/internal/logger"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler at user registration
func RegisterHandler(c *gin.Context) {
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

	resp.Code = 200
	resp.Msg = "success"
	c.JSON(http.StatusOK, resp)
}
