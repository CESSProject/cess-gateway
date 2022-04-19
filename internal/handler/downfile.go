package handler

import (
	"cess-httpservice/configs"
	"cess-httpservice/internal/chain"
	"cess-httpservice/internal/db"
	. "cess-httpservice/internal/logger"
	"cess-httpservice/internal/rpc"
	"cess-httpservice/tools"
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
	"storj.io/common/base58"
)

func DownfileHandler(c *gin.Context) {
	var resp = RespMsg{
		Code: http.StatusBadRequest,
		Msg:  "",
	}
	filename := c.Query("filename")
	if filename == "" {
		resp.Msg = "filename is empty"
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	// Determine if the user has uploaded the file
	key, err := tools.CalcMD5(filename)
	if err != nil {
		resp.Msg = "invalid filename"
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	db, err := db.GetDB()
	if err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	v, err := db.Get(key)
	if err != nil {
		if err.Error() == "leveldb: not found" {
			resp.Code = http.StatusNotFound
			resp.Msg = "This file has not been uploaded"
			c.JSON(http.StatusNotFound, resp)
			return
		} else {
			resp.Code = http.StatusInternalServerError
			resp.Msg = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
	}

	// local cache
	fdir := filepath.Join(configs.FileCacheDir, fmt.Sprintf("%v", string(key)))
	_, err = os.Stat(fdir)
	if err != nil {
		os.MkdirAll(fdir, os.ModeDir)
	}
	fpath := filepath.Join(fdir, filename)
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
		Err.Sugar().Errorf("%v", err)
		resp.Code = http.StatusInternalServerError
		resp.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	if string(filemetainfo.FileState) != "active" {
		resp.Code = http.StatusForbidden
		resp.Msg = "The file is in hot backup, please try again later."
		c.JSON(http.StatusForbidden, resp)
		return
	}

	// Download the file from the scheduler service
	err = downloadFromStorage(fpath, fid)
	if err != nil {
		Err.Sugar().Errorf("[%v] [%v] %v", fpath, fid, err)
		resp.Code = http.StatusInternalServerError
		resp.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	//c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filehash))
	c.Writer.Header().Add("inline", fmt.Sprintf("inline; filename=%v", filename))
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.File(fpath)
	defer os.Remove(fpath)
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
			if i == len(schds)-1 {
				return errors.New("All scheduler is offline")
			}
		} else {
			break
		}
	}

	var wantfile rpc.FileDownloadReq
	sp := sync.Pool{
		New: func() interface{} {
			return &rpc.ReqMsg{}
		},
	}
	wantfile.FileId = fmt.Sprintf("%v", fid)
	wantfile.WalletAddress = ""
	wantfile.Blocks = 1

	for {
		data, err := proto.Marshal(&wantfile)
		if err != nil {
			return err
		}
		req := sp.Get().(*rpc.ReqMsg)
		req.Method = configs.RpcMethod_ReadFile
		req.Service = configs.RpcService_Scheduler
		req.Body = data

		ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
		defer cancel()
		resp, err := client.Call(ctx, req)
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
		sp.Put(req)
	}
	return nil
}
