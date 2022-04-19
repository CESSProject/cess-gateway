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
	"io"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
	"storj.io/common/base58"
)

func UpfileHandler(c *gin.Context) {
	var resp = RespMsg{
		Code: http.StatusBadRequest,
		Msg:  "",
	}

	content_length := c.Request.ContentLength
	if content_length <= 0 {
		resp.Msg = "empty file"
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	file_p, err := c.FormFile("file")
	if err != nil {
		resp.Msg = "not upload file request"
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	spaceInfo, err := chain.GetUserSpaceInfo(configs.Confile.AccountAddr)
	if err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	if spaceInfo.Remaining_space.Uint64()*1024 < uint64(file_p.Size) {
		resp.Code = http.StatusForbidden
		resp.Msg = "Not enough free space"
		c.JSON(http.StatusForbidden, resp)
		return
	}

	file_c, _, _ := c.Request.FormFile("file")

	fnamemd5, err := tools.CalcMD5(file_p.Filename)
	if err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	//userpath := filepath.Join(configs.FileCacheDir, "test")
	userpath := filepath.Join(configs.FileCacheDir, fmt.Sprintf("%v", string(fnamemd5)))
	_, err = os.Stat(userpath)
	if err != nil {
		err = os.MkdirAll(userpath, os.ModeDir)
		if err != nil {
			resp.Code = http.StatusInternalServerError
			resp.Msg = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
	}

	fpath := filepath.Join(userpath, file_p.Filename)
	_, err = os.Stat(fpath)
	if err == nil {
		resp.Code = http.StatusForbidden
		resp.Msg = "duplicate file name"
		c.JSON(http.StatusForbidden, resp)
		return
	}

	resp.Code = http.StatusInternalServerError
	f, err := os.Create(fpath)
	if err != nil {
		Err.Sugar().Errorf("create file fail:%v\n", err)
		resp.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	defer f.Close()
	buf := make([]byte, 2*1024*1024)
	for {
		n, err := file_c.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			resp.Code = http.StatusGatewayTimeout
			resp.Msg = "upload failed due to network issues"
			c.JSON(http.StatusGatewayTimeout, resp)
			return
		}
		if n == 0 {
			continue
		}
		f.Write(buf[:n])
	}

	hash, err := tools.CalcFileHash2(f)
	if err != nil {
		Err.Sugar().Errorf("%v", err)
		resp.Code = http.StatusInternalServerError
		resp.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	db, err := db.GetDB()
	ok, _ := db.Has([]byte(hash))
	if ok {
		os.RemoveAll(fpath)
		resp.Code = http.StatusForbidden
		resp.Msg = "duplicate file hash"
		c.JSON(http.StatusForbidden, resp)
		return
	} else {
		db.Put([]byte(hash), []byte(fpath))
	}

	resp.Code = http.StatusOK
	resp.Msg = "success"
	c.JSON(http.StatusOK, resp)

	go uploadToStorage(fpath, string(fnamemd5))

	return
}

// Upload files to cess storage system
func uploadToStorage(fpath, fnamemd5 string) {
	time.Sleep(time.Second)
	defer func() {
		err := recover()
		if err != nil {
			Err.Sugar().Errorf("[panic]: %v", err)
		}
	}()
	file, err := os.Stat(fpath)
	if err != nil {
		Err.Sugar().Errorf("[%v] %v", fpath, err)
		return
	}

	filehash, err := tools.CalcFileHash(fpath)
	if err != nil {
		Err.Sugar().Errorf("[%v] %v", fpath, err)
		return
	}
	fileid, err := tools.GetGuid(int64(tools.RandomInRange(0, 1023)))
	if err != nil {
		Err.Sugar().Errorf("[%v] %v", fpath, err)
		return
	}
	var blockinfo rpc.FileUploadInfo
	blockinfo.Backups = "3"
	blockinfo.FileId = fmt.Sprintf("%v", fileid)
	blockinfo.BlockSize = 0
	blockinfo.FileHash = filehash

	blocksize := 2 * 1024 * 1024
	blocktotal := 0

	f, err := os.Open(fpath)
	if err != nil {
		Err.Sugar().Errorf("[%v] %v", fpath, err)
		return
	}
	defer f.Close()
	filebytes, err := ioutil.ReadAll(f)
	if err != nil {
		Err.Sugar().Errorf("[%v] %v", fpath, err)
		return
	}

	schds, err := chain.GetSchedulerInfo()
	if err != nil {
		Err.Sugar().Errorf("[%v] %v", fpath, err)
		return
	}

	var filesize int64 = 0
	if file.Size()/1024 == 0 {
		filesize = 1
	} else {
		filesize = file.Size() / 1024
	}

	err = chain.FileMetaInfoOnChain(
		configs.Confile.AccountSeed,
		configs.Confile.AccountAddr,
		file.Name(),
		fmt.Sprintf("%v", fileid),
		filehash,
		false,
		3,
		filesize,
		new(big.Int).SetUint64(0),
	)
	if err != nil {
		Err.Sugar().Errorf("[%v] %v", fpath, err)
		return
	}

	var client *rpc.Client
	for i, schd := range schds {
		wsURL := "ws://" + string(base58.Decode(string(schd.Ip)))
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		client, err = rpc.DialWebsocket(ctx, wsURL, "")
		//defer cancel()
		if err != nil {
			Err.Sugar().Errorf("[%v] %v", fpath, string(schds[i].Ip))
			if i == len(schds) {
				Err.Sugar().Errorf("[%v] All scheduler not working", len(schds))
				return
			}

		} else {
			break
		}
	}
	// sp := sync.Pool{
	// 	New: func() interface{} {
	// 		return &rpc.ReqMsg{}
	// 	},
	// }
	reqmsg := rpc.ReqMsg{}
	reqmsg.Method = configs.RpcMethod_WriteFile
	reqmsg.Service = configs.RpcService_Scheduler
	commit := func(num int, data []byte) error {
		blockinfo.BlockNum = int32(num) + 1
		blockinfo.Data = data
		info, err := proto.Marshal(&blockinfo)
		if err != nil {
			return errors.Wrap(err, "[Error]Serialization error, please upload again")
		}
		//reqmsg := sp.Get().(*rpc.ReqMsg)
		reqmsg.Body = info

		ctx, _ := context.WithTimeout(context.Background(), 90*time.Second)
		resp, err := client.Call(ctx, &reqmsg)
		//defer cancel()
		if err != nil {
			return errors.Wrap(err, "[Error]Failed to transfer file to scheduler,error")
		}

		var res rpc.RespBody
		err = proto.Unmarshal(resp.Body, &res)
		if err != nil {
			return errors.Wrap(err, "[Error]Error getting reply from schedule, transfer failed")
		}
		if res.Code != 0 {
			err = errors.New(res.Msg)
			return errors.Wrap(err, "[Error]Upload file fail!scheduler problem")
		}
		//sp.Put(reqmsg)
		return nil
	}
	blocks := len(filebytes) / blocksize
	if len(filebytes)%blocksize == 0 {
		blocktotal = blocks
	} else {
		blocktotal = blocks + 1
	}
	blockinfo.Blocks = int32(blocktotal)

	for i := 0; i < blocktotal; i++ {
		block := make([]byte, 0)
		if blocks != i {
			block = filebytes[i*blocksize : (i+1)*blocksize]
		} else {
			block = filebytes[i*blocksize:]
		}
		err = commit(i, block)
		if err != nil {
			Err.Sugar().Errorf("[%v] %v", fpath, err)
			return
		}
	}
	os.Remove(fpath)
	fmt.Printf("[Success] Storage file:%s successful", fpath)
	Out.Sugar().Infof("[Success] Storage file:%s successful", fpath)
	db, err := db.GetDB()
	if err != nil {
		Err.Sugar().Errorf("[%v][%v] %v", fpath, fileid, err)
		return
	}
	err = db.Put([]byte(fnamemd5), tools.Int64ToBytes(fileid))
	if err != nil {
		Err.Sugar().Errorf("[%v][%v] %v", fpath, fileid, err)
		return
	}
	err = db.Put(tools.Int64ToBytes(fileid), []byte(file.Name()))
	if err != nil {
		Err.Sugar().Errorf("[%v][%v] %v", fpath, fileid, err)
		return
	}
	fmt.Printf("[Success] DB record a file:%s successful", fpath)
	Out.Sugar().Infof("[Success] DB record a file:%s successful", fpath)
}
