package handler

import (
	"cess-gateway/configs"
	"cess-gateway/internal/chain"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Filestate_resp struct {
	Size  uint64
	State string
	Names []string
}

func FilestateHandler(c *gin.Context) {
	var resp = RespMsg{
		Code: http.StatusUnauthorized,
		Msg:  Status_401_token,
	}
	fid := c.Param("fid")
	fmt.Println("fid:", fid)
	if fid == "" {
		resp.Code = http.StatusBadRequest
		resp.Msg = Status_400_default
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	//query all file meta
	filestate, code, _ := chain.GetFileMetaInfoOnChain(fid)
	if code != configs.Code_200 {
		if code == configs.Code_404 {
			resp.Code = http.StatusNotFound
			resp.Msg = "Not found"
			c.JSON(http.StatusOK, resp)
			return
		}
		resp.Code = http.StatusInternalServerError
		resp.Msg = Status_500_chain
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	var fs Filestate_resp
	fs.Size = uint64(filestate.FileSize)
	fs.State = string(filestate.FileState)
	for _, v := range filestate.Names {
		var tmp string = string(v)
		fs.Names = append(fs.Names, tmp)
	}
	resp.Code = http.StatusOK
	resp.Msg = Status_200_default
	resp.Data = fs
	c.JSON(http.StatusOK, resp)
	return
}
