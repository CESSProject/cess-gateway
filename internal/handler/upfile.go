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
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	cesskeyring "github.com/CESSProject/go-keyring"
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
	if htoken == "" {
		Uld.Sugar().Infof("[%v] head missing token", c.ClientIP())
		c.JSON(http.StatusUnauthorized, resp)
		return
	}

	bytes, err := token.DecryptToken(htoken)
	if err != nil {
		Uld.Sugar().Infof("[%v] [%v] DecryptToken error", c.ClientIP(), htoken)
		c.JSON(http.StatusUnauthorized, resp)
		return
	}

	var usertoken token.TokenMsgType
	err = json.Unmarshal(bytes, &usertoken)
	if err != nil {
		Uld.Sugar().Infof("[%v] [%v] token format error", c.ClientIP(), htoken)
		c.JSON(http.StatusUnauthorized, resp)
		return
	}

	if time.Now().Unix() >= usertoken.ExpirationTime {
		Uld.Sugar().Infof("[%v] token expired", usertoken.Mailbox)
		resp.Msg = Status_401_expired
		c.JSON(http.StatusUnauthorized, resp)
		return
	}

	db, err := db.GetDB()
	if err != nil {
		Uld.Sugar().Infof("[%v] %v", usertoken.Mailbox, err)
		resp.Code = http.StatusInternalServerError
		resp.Msg = Status_500_db
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp.Code = http.StatusBadRequest
	resp.Msg = Status_400_default
	filename := c.Param("filename")
	if filename == "" {
		Uld.Sugar().Infof("[%v] no file name", usertoken.Mailbox)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	key, err := tools.CalcMD5(usertoken.Mailbox + url.QueryEscape(filename))
	if err != nil {
		Uld.Sugar().Infof("[%v] %v", usertoken.Mailbox, err)
		resp.Code = http.StatusInternalServerError
		resp.Msg = Status_500_unexpected
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	ok, err := db.Has(key)
	if err != nil {
		Uld.Sugar().Infof("[%v] %v", usertoken.Mailbox, err)
		resp.Code = http.StatusInternalServerError
		resp.Msg = Status_500_db
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	if ok {
		resp.Code = http.StatusForbidden
		resp.Msg = Status_403_dufilename
		c.JSON(http.StatusForbidden, resp)
		return
	}

	content_length := c.Request.ContentLength
	if content_length <= 0 {
		Uld.Sugar().Infof("[%v] contentLength <= 0", usertoken.Mailbox)
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	file_p, err := c.FormFile("file")
	if err != nil {
		Uld.Sugar().Infof("[%v] FormFile err", usertoken.Mailbox)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	spaceInfo, err := chain.GetUserSpaceInfo(configs.Confile.AccountSeed)
	if err != nil {
		Uld.Sugar().Infof("[%v] %v", usertoken.Mailbox, err)
		resp.Code = http.StatusInternalServerError
		resp.Msg = Status_500_chain
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	remainSpace := spaceInfo.Remaining_space.Uint64()

	if remainSpace < uint64(file_p.Size) {
		resp.Code = http.StatusForbidden
		resp.Msg = Status_403_NotEnoughSpace
		c.JSON(http.StatusForbidden, resp)
		return
	}

	file_c, _, err := c.Request.FormFile("file")
	if err != nil {
		Uld.Sugar().Infof("[%v] %v", usertoken.Mailbox, err)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	// server data
	resp.Code = http.StatusInternalServerError
	resp.Msg = Status_500_unexpected

	_, err = os.Stat(configs.FileCacheDir)
	if err != nil {
		err = os.MkdirAll(configs.FileCacheDir, os.ModeDir)
		if err != nil {
			Uld.Sugar().Infof("[%v] %v", usertoken.Mailbox, err)
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
	}

	fpath := filepath.Join(configs.FileCacheDir, url.QueryEscape(filename))
	_, err = os.Stat(fpath)
	if err == nil {
		Uld.Sugar().Infof("[%v] %v:%v", usertoken.Mailbox, Status_403_dufilename, fpath)
		resp.Code = http.StatusForbidden
		resp.Msg = Status_403_dufilename
		c.JSON(http.StatusForbidden, resp)
		return
	}

	f, err := os.Create(fpath)
	if err != nil {
		Uld.Sugar().Infof("[%v] %v", usertoken.Mailbox, err)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
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
	f.Close()

	//Calc file id
	hash, err := calcFileHashByChunks(fpath, configs.SIZE_1GB)
	if err != nil {
		Uld.Sugar().Infof("[%v] %v", usertoken.Mailbox, err)
		resp.Msg = Status_500_unexpected
		c.JSON(http.StatusInternalServerError, resp)
	}
	fileid := "cess" + hash
	newpath := filepath.Join(configs.FileCacheDir, fileid)
	os.Rename(fpath, newpath)

	key_fid := usertoken.Mailbox + fileid

	txhash, _, err := chain.UploadDeclaration(configs.Confile.AccountSeed, fileid, filename)
	if txhash == "" {
		Uld.Sugar().Infof("[%v] %v", usertoken.Mailbox, err)
		resp.Msg = Status_500_db
		c.JSON(http.StatusInternalServerError, resp)
	}

	err = db.Put([]byte(key), []byte(key_fid))
	if err != nil {
		Uld.Sugar().Infof("[%v] %v", usertoken.Mailbox, err)
		resp.Msg = Status_500_db
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	err = db.Put([]byte(key_fid), []byte(key))
	if err != nil {
		Uld.Sugar().Infof("[%v] %v", usertoken.Mailbox, err)
		resp.Msg = Status_500_db
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	go task_StoreFile(newpath, usertoken.Mailbox, fileid, filename)
	resp.Code = http.StatusOK
	resp.Msg = Status_200_default
	resp.Data = fmt.Sprintf("%v", fileid)
	c.JSON(http.StatusOK, resp)
	return
}

func task_StoreFile(fpath, mailbox, fid, fname string) {
	defer func() {
		if err := recover(); err != nil {
			Err.Sugar().Errorf("%v", err)
		}
	}()
	var channel_1 = make(chan uint8, 1)
	Uld.Sugar().Infof("[%v] Start the file backup management process", fid)
	go uploadToStorage(channel_1, fpath, mailbox, fid, fname)
	for {
		select {
		case result := <-channel_1:
			if result == 1 {
				go uploadToStorage(channel_1, fpath, mailbox, fid, fname)
			}
			if result == 2 {
				Uld.Sugar().Infof("[%v] File save successfully", fid)
				return
			}
			if result == 3 {
				Uld.Sugar().Infof("[%v] File save failed", fid)
				return
			}
		}
	}
}

// Upload files to cess storage system
func uploadToStorage(ch chan uint8, fpath, mailbox, fid, fname string) {
	defer func() {
		err := recover()
		if err != nil {
			ch <- 1
			Uld.Sugar().Infof("[panic]: [%v] [%v] %v", mailbox, fpath, err)
		}
	}()
	fstat, err := os.Stat(fpath)
	if err != nil {
		ch <- 3
		Uld.Sugar().Infof("[%v] [%v] %v", mailbox, fpath, err)
		return
	}

	var authreq rpc.AuthReq
	authreq.FileId = fid
	authreq.FileName = fname
	authreq.FileSize = uint64(fstat.Size())
	authreq.BlockTotal = uint32(fstat.Size() / configs.RpcBuffer)
	if fstat.Size()%configs.RpcBuffer != 0 {
		authreq.BlockTotal += 1
	}
	authreq.PublicKey, err = chain.GetPubkeyFromPrk(configs.Confile.AccountSeed)
	if err != nil {
		ch <- 1
		Uld.Sugar().Infof("[%v] [%v] %v", mailbox, fpath, err)
		return
	}

	authreq.Msg = []byte(tools.GetRandomcode(16))
	kr, _ := cesskeyring.FromURI(configs.Confile.AccountSeed, cesskeyring.NetSubstrate{})
	// sign message
	sign, err := kr.Sign(kr.SigningContext(authreq.Msg))
	if err != nil {
		ch <- 1
		Uld.Sugar().Infof("[%v] [%v] %v", mailbox, fpath, err)
		return
	}
	authreq.Sign = sign[:]

	// get all scheduler
	schds, err := chain.GetSchedulerInfo()
	if err != nil {
		ch <- 1
		Uld.Sugar().Infof("[%v] [%v] %v", mailbox, fpath, err)
		return
	}
	var client *rpc.Client
	for i, schd := range schds {
		wsURL := "ws://" + string(base58.Decode(string(schd.Ip)))
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		client, err = rpc.DialWebsocket(ctx, wsURL, "")
		if err != nil {
			Uld.Sugar().Infof("[%v] [%v] [%v] %v", mailbox, fpath, wsURL, err)
			if (i + 1) == len(schds) {
				ch <- 1
				Uld.Sugar().Infof("[%v] [%v] All scheduler not working", mailbox, fpath)
				return
			}
		} else {
			break
		}
	}

	bob, err := proto.Marshal(&authreq)
	if err != nil {
		ch <- 1
		Uld.Sugar().Infof("[%v] [%v] %v", mailbox, fpath, err)
		return
	}

	data, code, err := WriteData2(client, configs.RpcService_Scheduler, configs.RpcMethod_auth, bob)
	if err != nil {
		ch <- 1
		Uld.Sugar().Infof("[%v] [%v] %v", mailbox, fpath, err)
		return
	}

	if code == 201 {
		ch <- 2
		return
	}

	if code != 200 {
		ch <- 1
		Uld.Sugar().Infof("[%v] [%v] %v", mailbox, fpath, code)
		return
	}

	var n int
	var filereq rpc.FileUploadReq
	var buf = make([]byte, configs.RpcBuffer)
	f, err := os.OpenFile(fpath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		ch <- 1
		Uld.Sugar().Infof("[%v] [%v] %v", mailbox, fpath, code)
		return
	}
	filereq.Auth = data
	for i := 0; i < int(authreq.BlockTotal); i++ {
		filereq.BlockIndex = uint32(i + 1)
		f.Seek(int64(i*configs.RpcBuffer), 0)
		n, _ = f.Read(buf)
		filereq.FileData = buf[:n]

		bob, err := proto.Marshal(&filereq)
		if err != nil {
			ch <- 1
			Uld.Sugar().Infof("[%v] [%v] %v", mailbox, fpath, err)
			return
		}

		_, _, err = WriteData2(client, configs.RpcService_Scheduler, configs.RpcMethod_WriteFile, bob)
		if err != nil {
			ch <- 1
			Uld.Sugar().Infof("[%v] [%v] %v", mailbox, fpath, err)
			return
		}
	}
	ch <- 2
	Uld.Sugar().Infof("[Success] Storage file:%s", fpath)
}

func calcFileHashByChunks(fpath string, chunksize int64) (string, error) {
	if chunksize <= 0 {
		return "", errors.New("Invalid chunk size")
	}
	fstat, err := os.Stat(fpath)
	if err != nil {
		return "", err
	}
	chunkNum := fstat.Size() / chunksize
	if fstat.Size()%chunksize != 0 {
		chunkNum++
	}
	var n int
	var chunkhash, allhash, filehash string
	var buf = make([]byte, chunksize)
	f, err := os.OpenFile(fpath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return "", err
	}
	defer f.Close()
	for i := int64(0); i < chunkNum; i++ {
		f.Seek(i*chunksize, 0)
		n, err = f.Read(buf)
		if err != nil && err != io.EOF {
			return "", err
		}
		chunkhash, err = tools.CalcHash(buf[:n])
		if err != nil {
			return "", err
		}
		allhash += chunkhash
	}
	filehash, err = tools.CalcHash([]byte(allhash))
	if err != nil {
		return "", err
	}
	return filehash, nil
}

func WriteData2(cli *rpc.Client, service, method string, body []byte) ([]byte, int32, error) {
	req := &rpc.ReqMsg{
		Service: service,
		Method:  method,
		Body:    body,
	}
	ctx, _ := context.WithTimeout(context.Background(), 90*time.Second)
	resp, err := cli.Call(ctx, req)
	if err != nil {
		return nil, 0, errors.Wrap(err, "Call err:")
	}

	var b rpc.RespBody
	err = proto.Unmarshal(resp.Body, &b)
	if err != nil {
		return nil, 0, errors.Wrap(err, "Unmarshal:")
	}
	return b.Data, b.Code, err
}
