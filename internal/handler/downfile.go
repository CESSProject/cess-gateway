package handler

import (
	"cess-gateway/configs"
	"cess-gateway/internal/chain"
	"cess-gateway/internal/fileHandling"
	. "cess-gateway/internal/logger"
	"cess-gateway/internal/rpc"
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
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
	fmeta, err := chain.GetFileMetaInfoOnChain(fid)
	if err != nil {
		Err.Sugar().Errorf("[%v] %v", c.ClientIP(), err)
		if err.Error() == chain.ERR_Empty {
			resp.Code = http.StatusNotFound
			resp.Msg = Status_400_NotUploaded
			c.JSON(http.StatusNotFound, resp)
			return
		}
		resp.Code = http.StatusInternalServerError
		resp.Msg = Status_500_unexpected
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	if string(fmeta.FileState) != "active" {
		Err.Sugar().Errorf("[%v] file state is not active", c.ClientIP())
		resp.Code = http.StatusForbidden
		resp.Msg = Status_403_hotbackup
		c.JSON(http.StatusForbidden, resp)
		return
	}

	for i := 0; i < len(fmeta.ChunkInfo); i++ {
		// Download the file from the scheduler service
		fname := filepath.Join(configs.FileCacheDir, string(fmeta.ChunkInfo[i].ChunkId))
		err = downloadFromStorage(fname, string(fmeta.ChunkInfo[i].MinerIp))
		if err != nil {
			Err.Sugar().Errorf("[%v] Error downloading %drd shard", c.ClientIP(), i)
		}
	}
	r := len(fmeta.ChunkInfo) / 3
	d := len(fmeta.ChunkInfo) - r
	err = fileHandling.ReedSolomon_Restore(configs.FileCacheDir, fid, d, r)
	if err != nil {
		Err.Sugar().Errorf("[%v] ReedSolomon_Restore: %v", c.ClientIP(), err)
		resp.Code = http.StatusInternalServerError
		resp.Msg = Status_500_unexpected
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	//c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filehash))
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%v", fid))
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.File(fpath)
	return
}

// Download files from cess storage service
func downloadFromStorage(fpath string, mip string) error {
	file, err := os.OpenFile(fpath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	var client *rpc.Client

	wsURL := "ws://" + string(base58.Decode(mip))

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	client, err = rpc.DialWebsocket(ctx, wsURL, "")
	if err != nil {
		return err
	}

	var wantfile rpc.FileDownloadReq
	fname := filepath.Base(fpath)

	wantfile.FileId = fmt.Sprintf("%v", fname)
	wantfile.BlockIndex = 1

	reqmsg := rpc.ReqMsg{}
	reqmsg.Method = configs.RpcMethod_ReadFile
	reqmsg.Service = configs.RpcService_Miner
	for {
		data, err := proto.Marshal(&wantfile)
		if err != nil {
			return err
		}
		reqmsg.Body = data
		ctx, _ := context.WithTimeout(context.Background(), 90*time.Second)
		resp, err := client.Call(ctx, &reqmsg)
		if err != nil {
			return err
		}

		var respbody rpc.RespBody
		err = proto.Unmarshal(resp.Body, &respbody)
		if err != nil || respbody.Code != 200 {
			return errors.Wrap(err, "[Error]Download file from CESS reply message"+respbody.Msg+",error")
		}
		var blockData rpc.FileDownloadInfo
		err = proto.Unmarshal(respbody.Data, &blockData)
		if err != nil {
			return errors.Wrap(err, "[Error]Download file from CESS error")
		}

		_, err = file.Write(blockData.Data)
		if err != nil {
			return err
		}

		if blockData.BlockIndex == blockData.BlockTotal {
			break
		}
		wantfile.BlockIndex++
	}
	return nil
}
