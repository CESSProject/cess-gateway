package handler

import (
	"cess-gateway/configs"
	"cess-gateway/internal/chain"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UserStateHandler(c *gin.Context) {
	var resp = RespMsg{
		Code: http.StatusBadRequest,
		Msg:  "",
	}

	spaceDetailsList, err := chain.GetSpaceDetailsInfo(configs.Confile.AccountAddr)
	if err != nil {
		resp.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	spaceInfo, err := chain.GetUserSpaceInfo(configs.Confile.AccountAddr)
	if err != nil {
		resp.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	lastestBlockHeight, err := chain.GetLastestBlockHeight()
	if err != nil {
		resp.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	var userdata UserStateMsg
	var spacelist = make([]SpaceDetailsMsg, 0)
	userdata.TotalSpace = fmt.Sprintf("%v", spaceInfo.Purchased_space.Uint64()) + " kb"
	userdata.UsedSpace = fmt.Sprintf("%v", spaceInfo.Used_space.Uint64()) + " kb"
	userdata.FreeSpace = fmt.Sprintf("%v", spaceInfo.Remaining_space.Uint64()) + " kb"
	for i := 0; i < len(spaceDetailsList); i++ {
		spacelist[i].Size = spaceDetailsList[i].Size.Uint64()
		if uint32(spaceDetailsList[i].Deadline) > lastestBlockHeight {
			spacelist[i].Deadline = (uint32(spaceDetailsList[i].Deadline) - lastestBlockHeight) * 3
		} else {
			spacelist[i].Deadline = 0
		}
	}
	userdata.SpaceDetails = spacelist
	resp.Code = http.StatusOK
	resp.Msg = "success"
	resp.Data = userdata
	c.JSON(http.StatusOK, resp)
}
