package handler

import (
	"bufio"
	"cess-gateway/configs"
	. "cess-gateway/internal/logger"
	"cess-gateway/internal/token"
	"cess-gateway/tools"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/btcsuite/btcutil/base58"
	"github.com/gin-gonic/gin"
)

func FilelistHandler(c *gin.Context) {
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

	// Parameters
	resp.Code = http.StatusBadRequest
	resp.Msg = Status_400_default
	var page, size, strartIndex = 0, 0, 0
	var defaultPage, defaultSize = true, true
	sizes := c.Query("size")
	pages := c.Query("page")
	if pages != "" {
		page, err = strconv.Atoi(pages)
		if err != nil {
			Err.Sugar().Errorf("[%v] [%v] filename is empty", c.ClientIP(), usertoken.Mailbox)
			c.JSON(http.StatusBadRequest, resp)
			return
		}
		defaultPage = false
	}
	if sizes != "" {
		size, err = strconv.Atoi(sizes)
		if err != nil {
			Err.Sugar().Errorf("[%v] [%v] filename is empty", c.ClientIP(), usertoken.Mailbox)
			c.JSON(http.StatusBadRequest, resp)
			return
		}
		defaultSize = false
	}
	resp.Code = http.StatusInternalServerError
	resp.Msg = Status_500_unexpected
	fs, err := tools.WalkDir(filepath.Join(configs.FileCacheDir, fmt.Sprintf("%v", usertoken.UserId), configs.FilRecordsDir))
	if err != nil {
		Err.Sugar().Errorf("[%v] [%v] %v", c.ClientIP(), usertoken.Mailbox, err)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	sort.Strings(fs)
	if len(fs) > 0 {
		Err.Sugar().Errorf("[%v] [%v] [%v] not found records file", c.ClientIP(), usertoken.Mailbox, filepath.Join(configs.FileCacheDir, fmt.Sprintf("%v", usertoken.UserId)))
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	if defaultPage {
		if defaultSize {
			size = 30
		} else {
			if size > 1000 {
				size = 1000
			}
		}
		var fnamelist = make([]string, size)
		file, err := os.Open(filepath.Join(configs.FileCacheDir, fmt.Sprintf("%v", usertoken.UserId), configs.FilRecordsDir, fs[len(fs)-1]))
		if err != nil {
			Err.Sugar().Errorf("[%v] [%v] %v", c.ClientIP(), usertoken.Mailbox, err)
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
		defer file.Close()
		buffer := bufio.NewReader(file)
		for {
			ctx, _, err := buffer.ReadLine()
			if err != nil {
				break
			}
			if strings.TrimSpace(string(ctx)) == "" {
				continue
			}
			fnamelist = append(fnamelist, string(ctx))
		}
		if len(fnamelist) < size && len(fs) > 1 {
			file, err := os.Open(filepath.Join(configs.FileCacheDir, fmt.Sprintf("%v", usertoken.UserId), configs.FilRecordsDir, fs[len(fs)-2]))
			if err != nil {
				Err.Sugar().Errorf("[%v] [%v] %v", c.ClientIP(), usertoken.Mailbox, err)
				c.JSON(http.StatusInternalServerError, resp)
				return
			}
			defer file.Close()
			var fnamelist_pre = make([]string, 1000)
			buffer := bufio.NewReader(file)
			for {
				ctx, _, err := buffer.ReadLine()
				if err != nil {
					break
				}
				if strings.TrimSpace(string(ctx)) == "" {
					continue
				}
				fnamelist_pre = append(fnamelist_pre, string(ctx))
			}
			if (size - len(fnamelist)) > len(fnamelist_pre) {
				fnamelist = append(fnamelist, fnamelist_pre...)
			} else {
				fnamelist = append(fnamelist, fnamelist_pre[(len(fnamelist_pre)+len(fnamelist)-size):]...)
			}
		}
		var data_names = make([]string, 0)
		if len(fnamelist) <= size {
			for i := range fnamelist {
				data_names = append(data_names, string(base58.Decode(fnamelist[i])))
			}
		} else {
			for i := 0; i < size; i++ {
				data_names = append(data_names, string(base58.Decode(fnamelist[len(fnamelist)-size+i])))
			}
		}
		resp.Code = http.StatusOK
		resp.Msg = "success"
		resp.Data = data_names
		return
	} else {
		strartIndex = page * 30
		strartIndex = strartIndex
	}

	// strartIndex :=
	resp.Code = http.StatusOK
	resp.Msg = "success"
	//resp.Data = fnames
	c.JSON(http.StatusOK, resp)
}
