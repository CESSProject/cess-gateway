package handler

import (
	"cess-gateway/configs"
	"cess-gateway/internal/chain"
	"cess-gateway/internal/db"
	. "cess-gateway/internal/logger"
	"cess-gateway/internal/rpc"
	"cess-gateway/internal/token"
	"cess-gateway/tools"
	"context"
	"encoding/json"
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

	// token
	htoken := c.Request.Header.Get("Authorization")
	if htoken == "" {
		Err.Sugar().Errorf("[%v] head missing token", c.ClientIP())
		c.JSON(http.StatusUnauthorized, resp)
		return
	}

	bytes, err := token.DecryptToken(htoken)
	if err != nil {
		Err.Sugar().Errorf("[%v] [%v] DecryptToken error", c.ClientIP(), htoken)
		c.JSON(http.StatusUnauthorized, resp)
		return
	}

	var usertoken token.TokenMsgType
	err = json.Unmarshal(bytes, &usertoken)
	if err != nil {
		Err.Sugar().Errorf("[%v] [%v] token format error", c.ClientIP(), htoken)
		c.JSON(http.StatusUnauthorized, resp)
		return
	}

	if time.Now().Unix() >= usertoken.ExpirationTime {
		Err.Sugar().Errorf("[%v] [%v] token expired", c.ClientIP(), usertoken.Mailbox)
		resp.Msg = Status_401_expired
		c.JSON(http.StatusUnauthorized, resp)
		return
	}

	filename := c.Query("filename")
	if filename == "" {
		Err.Sugar().Errorf("[%v] [%v] filename is empty", c.ClientIP(), usertoken.Mailbox)
		resp.Code = http.StatusBadRequest
		resp.Msg = Status_400_default
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	// server
	resp.Code = http.StatusInternalServerError
	resp.Msg = Status_500_unexpected
	key, err := tools.CalcMD5(usertoken.Mailbox + filename)
	if err != nil {
		Err.Sugar().Errorf("[%v] [%v] %v", c.ClientIP(), usertoken.Mailbox, err)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	db, err := db.GetDB()
	if err != nil {
		Err.Sugar().Errorf("[%v] [%v] %v", c.ClientIP(), usertoken.Mailbox, err)
		resp.Msg = Status_500_db
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	v, err := db.Get(key)
	if err != nil {
		if err.Error() == "leveldb: not found" {
			resp.Code = http.StatusBadRequest
			resp.Msg = Status_400_NotUploaded
			c.JSON(http.StatusBadRequest, resp)
			return
		} else {
			Err.Sugar().Errorf("[%v] [%v] %v", c.ClientIP(), usertoken.Mailbox, err)
			resp.Msg = Status_500_db
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
	}

	// local cache
	fdir := filepath.Join(configs.FileCacheDir, fmt.Sprintf("%v", usertoken.UserId))
	_, err = os.Stat(fdir)
	if err != nil {
		os.MkdirAll(fdir, os.ModeDir)
	}
	fpath := filepath.Join(fdir, filename)
	defer os.Remove(fpath)
	_, err = os.Stat(fpath)
	if err == nil {
		c.Writer.Header().Add("inline", fmt.Sprintf("inline; filename=%v", filename))
		c.Writer.Header().Add("Content-Type", "application/octet-stream")
		c.File(fpath)
		return
	}
	fid := tools.BytesToInt64(v)
	// file meta info
	filemetainfo, err := chain.GetFileMetaInfo(fid)
	if err != nil {
		Err.Sugar().Errorf("[%v] [%v] %v", c.ClientIP(), usertoken.Mailbox, err)
		resp.Msg = Status_500_chain
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	if string(filemetainfo.FileState) != "active" {
		resp.Code = http.StatusForbidden
		resp.Msg = Status_403_hotbackup
		c.JSON(http.StatusForbidden, resp)
		return
	}

	// Download the file from the scheduler service
	err = downloadFromStorage(fpath, fid)
	if err != nil {
		Err.Sugar().Errorf("[%v] [%v] %v", c.ClientIP(), usertoken.Mailbox, err)
		resp.Code = http.StatusInternalServerError
		resp.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	//c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filehash))
	c.Writer.Header().Add("inline", fmt.Sprintf("inline; filename=%v", filename))
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.File(fpath)
	return
}

// Download files from cess storage service
func downloadFromStorage(fpath string, fid int64) error {
	file, err := os.OpenFile(fpath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	schds, err := chain.GetSchedulerInfo()
	if err != nil {
		return err
	}

	var client *rpc.Client
	for i, schd := range schds {
		wsURL := "ws://" + string(base58.Decode(string(schd.Ip)))
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		client, err = rpc.DialWebsocket(ctx, wsURL, "")
		defer cancel()
		if err != nil {
			Err.Sugar().Errorf("[%v] %v", fpath, string(schds[i].Ip))
			if (i + 1) == len(schds) {
				return errors.New("All scheduler is offline")
			}
		} else {
			break
		}
	}

	var wantfile rpc.FileDownloadReq

	wantfile.FileId = fmt.Sprintf("%v", fid)
	wantfile.WalletAddress = ""
	wantfile.Blocks = 1

	reqmsg := rpc.ReqMsg{}
	reqmsg.Method = configs.RpcMethod_ReadFile
	reqmsg.Service = configs.RpcService_Scheduler
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
		if err != nil || respbody.Code != 0 {
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

		if blockData.Blocks == blockData.BlockNum {
			break
		}
		wantfile.Blocks++
	}
	return nil
}
