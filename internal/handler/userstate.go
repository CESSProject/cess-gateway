package handler

import (
	"cess-httpservice/configs"
	"cess-httpservice/internal/chain"
	"cess-httpservice/internal/token"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UserStateHandler(c *gin.Context) {
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
	collaterals, err := chain.GetUserInfo(token_de.Walletaddr)
	if err != nil {
		resp.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	spaceInfo, err := chain.GetUserSpaceInfo(token_de.Walletaddr)
	if err != nil {
		resp.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	var userdata UserStateMsg
	temp1, _ := new(big.Int).SetString(configs.CessTokenAccuracy, 10)
	temp2, _ := new(big.Int).SetString(collaterals.Collaterals.String(), 10)
	temp2.Div(temp2, temp1)
	userdata.UserId = token_de.Userid
	userdata.Walletaddr = token_de.Walletaddr
	userdata.Deposit = fmt.Sprintf("%v", temp2.Uint64()) + " CESS"
	userdata.TotalSpace = fmt.Sprintf("%v", spaceInfo.Purchased_space.Uint64()) + " kb"
	userdata.UsedSpace = fmt.Sprintf("%v", spaceInfo.Used_space.Uint64()) + " kb"
	userdata.FreeSpace = fmt.Sprintf("%v", spaceInfo.Remaining_space.Uint64()) + " kb"
	resp.Code = http.StatusOK
	resp.Msg = "success"
	resp.Data = userdata
	c.JSON(http.StatusOK, resp)
}
