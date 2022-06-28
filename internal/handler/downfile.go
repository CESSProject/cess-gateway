package handler

import (
	"cess-gateway/configs"
	"cess-gateway/internal/chain"
	. "cess-gateway/internal/logger"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"storj.io/common/base58"
)

func DownfileHandler(c *gin.Context) {
	var resp = RespMsg{
		Code: http.StatusUnauthorized,
		Msg:  Status_401_token,
	}

	fid := c.Param("fid")
	if fid == "" {
		Err.Sugar().Errorf("[%v] fid is empty", c.ClientIP())
		resp.Code = http.StatusBadRequest
		resp.Msg = Status_400_default
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	// local cache
	fpath := filepath.Join(configs.FileCacheDir, fid)
	_, err := os.Stat(fpath)
	if err == nil {
		//c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filehash))
		c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%v", fid))
		c.Writer.Header().Add("Content-Type", "application/octet-stream")
		c.File(fpath)
		return
	}

	// file meta info
	count := 0
	code := configs.Code_404
	fmeta := chain.FileMetaInfo{}
	for code != configs.Code_200 {
		fmeta, code, err = chain.GetFileMetaInfoOnChain(fid)
		if count > 3 && code != configs.Code_200 {
			Err.Sugar().Errorf("[%v] %v", c.ClientIP(), err)
			resp.Code = http.StatusInternalServerError
			resp.Msg = Status_500_unexpected
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
		if code != configs.Code_200 {
			time.Sleep(time.Second * 3)
		} else {
			if string(fmeta.FileState) != "active" {
				Err.Sugar().Errorf("[%v] file state is not active", c.ClientIP())
				resp.Code = http.StatusBadRequest
				resp.Msg = Status_403_default
				c.JSON(http.StatusForbidden, resp)
				return
			}
		}
		count++
	}
	dstip := "http://" + string(base58.Decode(string(fmeta.MinerIp))) + "/" + fid
	c.Redirect(http.StatusTemporaryRedirect, dstip)
}
