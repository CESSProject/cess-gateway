package handler

import (
	"cess-httpservice/internal/token"
	"encoding/json"
	"fmt"
	"net/http"

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

	//c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filehash))
	c.Writer.Header().Add("inline", fmt.Sprintf("inline; filename=%v", filename))
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.File("/usr/local/cess/cache/test/goo.mod")
}
