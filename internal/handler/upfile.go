package handler

import (
	"cess-httpservice/configs"
	"cess-httpservice/internal/chain"
	"cess-httpservice/internal/db"
	. "cess-httpservice/internal/logger"
	"cess-httpservice/internal/rpc"
	"cess-httpservice/internal/token"
	"cess-httpservice/tools"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
	"storj.io/common/base58"
)

func UpfileHandler(c *gin.Context) {
	var resp = RespMsg{
		Code: http.StatusUnauthorized,
		Msg:  Status_401_token,
	}
	// token
	htoken := c.Request.Header.Get("Authorization")
	fmt.Println(htoken)
	if htoken == "" {
		Err.Sugar().Errorf("[%v] head missing token", c.ClientIP())
		c.JSON(http.StatusUnauthorized, resp)
		return
	}
	var usertoken token.TokenMsgType
	err := json.Unmarshal([]byte(htoken), &usertoken)
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

	// client data
	resp.Code = http.StatusBadRequest
	resp.Msg = Status_400_default
	content_length := c.Request.ContentLength
	if content_length <= 0 {
		Err.Sugar().Errorf("[%v] [%v] contentLength <= 0", c.ClientIP(), usertoken.Mailbox)
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	file_p, err := c.FormFile("file")
	if err != nil {
		Err.Sugar().Errorf("[%v] [%v] FormFile err", c.ClientIP(), usertoken.Mailbox)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	spaceInfo, err := chain.GetUserSpaceInfo(configs.Confile.AccountAddr)
	if err != nil {
		Err.Sugar().Errorf("[%v] [%v] %v", c.ClientIP(), usertoken.Mailbox, err)
		resp.Code = http.StatusInternalServerError
		resp.Msg = Status_500_chain
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	if spaceInfo.Remaining_space.Uint64()*1024 < uint64(file_p.Size) {
		resp.Code = http.StatusForbidden
		resp.Msg = Status_403_expired
		c.JSON(http.StatusForbidden, resp)
		return
	}

	file_c, _, err := c.Request.FormFile("file")
	if err != nil {
		Err.Sugar().Errorf("[%v] [%v] %v", c.ClientIP(), usertoken.Mailbox, err)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	// server data
	resp.Code = http.StatusInternalServerError
	resp.Msg = Status_500_unexpected
	userpath := filepath.Join(configs.FileCacheDir, fmt.Sprintf("%v", usertoken.UserId))
	_, err = os.Stat(userpath)
	if err != nil {
		err = os.MkdirAll(userpath, os.ModeDir)
		if err != nil {
			Err.Sugar().Errorf("[%v] [%v] %v", c.ClientIP(), usertoken.Mailbox, err)
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
	}

	fpath := filepath.Join(userpath, file_p.Filename)
	_, err = os.Stat(fpath)
	if err == nil {
		Err.Sugar().Errorf("[%v] [%v] %v:%v", c.ClientIP(), usertoken.Mailbox, Status_403_dufilename, fpath)
		resp.Code = http.StatusForbidden
		resp.Msg = Status_403_dufilename
		c.JSON(http.StatusForbidden, resp)
		return
	}

	f, err := os.Create(fpath)
	if err != nil {
		Err.Sugar().Errorf("[%v] [%v] %v", c.ClientIP(), usertoken.Mailbox, err)
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

	fileid, err := tools.GetGuid(int64(tools.RandomInRange(0, 1023)))
	if err != nil {
		Err.Sugar().Errorf("[%v] [%v] %v", c.ClientIP(), usertoken.Mailbox, err)
		return
	}
	db, err := db.GetDB()
	if err != nil {
		Err.Sugar().Errorf("[%v] [%v] %v", c.ClientIP(), usertoken.Mailbox, err)
		resp.Msg = Status_500_db
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	fkey, err := tools.CalcMD5(usertoken.Mailbox + file_p.Filename)
	if err != nil {
		Err.Sugar().Errorf("[%v] [%v] %v", c.ClientIP(), usertoken.Mailbox, err)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	err = db.Put([]byte(fkey), tools.Int64ToBytes(fileid))
	if err != nil {
		Err.Sugar().Errorf("[%v] [%v] %v", c.ClientIP(), usertoken.Mailbox, err)
		resp.Msg = Status_500_db
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	err = db.Put(tools.Int64ToBytes(fileid), []byte(file_p.Filename))
	if err != nil {
		Err.Sugar().Errorf("[%v] [%v] %v", c.ClientIP(), usertoken.Mailbox, err)
		resp.Msg = Status_500_db
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp.Code = http.StatusOK
	resp.Msg = Status_200_default
	c.JSON(http.StatusOK, resp)

	go uploadToStorage(fpath, usertoken.Mailbox, fileid)

	return
}

// Upload files to cess storage system
func uploadToStorage(fpath, mailbox string, fid int64) {
	time.Sleep(time.Second)
	defer func() {
		err := recover()
		if err != nil {
			Err.Sugar().Errorf("[panic]: [%v] [%v] %v", mailbox, fpath, err)
		}
	}()
	file, err := os.Stat(fpath)
	if err != nil {
		Err.Sugar().Errorf("[%v] [%v] %v", mailbox, fpath, err)
		return
	}

	filehash, err := tools.CalcFileHash(fpath)
	if err != nil {
		Err.Sugar().Errorf("[%v] [%v] %v", mailbox, fpath, err)
		return
	}

	var blockinfo rpc.FileUploadInfo
	blockinfo.Backups = "3"
	blockinfo.FileId = strconv.FormatInt(fid, 10)
	blockinfo.BlockSize = 0
	blockinfo.FileHash = filehash

	blocksize := 2 * 1024 * 1024
	blocktotal := 0

	f, err := os.Open(fpath)
	if err != nil {
		Err.Sugar().Errorf("[%v] [%v] %v", mailbox, fpath, err)
		return
	}
	defer f.Close()
	filebytes, err := ioutil.ReadAll(f)
	if err != nil {
		Err.Sugar().Errorf("[%v] [%v] %v", mailbox, fpath, err)
		return
	}

	schds, err := chain.GetSchedulerInfo()
	if err != nil {
		Err.Sugar().Errorf("[%v] [%v] %v", mailbox, fpath, err)
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
		strconv.FormatInt(fid, 10),
		filehash,
		false,
		3,
		filesize,
		new(big.Int).SetUint64(0),
	)
	if err != nil {
		Err.Sugar().Errorf("[%v] [%v] %v", mailbox, fpath, err)
		return
	}

	var client *rpc.Client
	for i, schd := range schds {
		wsURL := "ws://" + string(base58.Decode(string(schd.Ip)))
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		client, err = rpc.DialWebsocket(ctx, wsURL, "")
		if err != nil {
			Err.Sugar().Errorf("[%v] [%v] [%v] %v", mailbox, fpath, wsURL, err)
			if i == len(schds) {
				Err.Sugar().Errorf("[%v] [%v] All scheduler not working", mailbox, fpath)
				return
			}
		} else {
			break
		}
	}
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
		reqmsg.Body = info

		ctx, _ := context.WithTimeout(context.Background(), 90*time.Second)
		resp, err := client.Call(ctx, &reqmsg)
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
	Out.Sugar().Infof("[Success] Storage file:%s", fpath)
}
