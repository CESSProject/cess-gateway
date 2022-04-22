package handler

import (
	"cess-httpservice/configs"
	"cess-httpservice/internal/communication"
	"cess-httpservice/internal/db"
	. "cess-httpservice/internal/logger"
	"cess-httpservice/internal/token"
	"cess-httpservice/tools"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// It is used to authorize users
func GrantTokenHandler(c *gin.Context) {
	var resp = RespMsg{
		Code: http.StatusBadRequest,
		Msg:  Status_400_default,
	}
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		Err.Sugar().Errorf("%v,%v", c.ClientIP(), err)
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	var reqmsg ReqGrantMsg
	err = json.Unmarshal(body, &reqmsg)
	if err != nil {
		Err.Sugar().Errorf("%v,%v", c.ClientIP(), err)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	// Check if the email format is correct
	if !communication.VerifyMailboxFormat(reqmsg.Mailbox) {
		Err.Sugar().Errorf("%v,%v", c.ClientIP(), err)
		resp.Msg = Status_400_EmailFormat
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	resp.Code = http.StatusInternalServerError
	db, err := db.GetDB()
	if err != nil {
		Err.Sugar().Errorf("[%v] [%v] %v", c.ClientIP(), reqmsg, err)
		resp.Msg = Status_500_db
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	bytes, err := db.Get([]byte(reqmsg.Mailbox))
	if err != nil {
		if err.Error() == "leveldb: not found" {
			captcha := tools.RandomInRange(100000, 999999)
			v := fmt.Sprintf("%v", captcha) + "#" + fmt.Sprintf("%v", time.Now().Add(time.Minute*10).Unix())
			err = db.Put([]byte(reqmsg.Mailbox), []byte(v))
			if err != nil {
				Err.Sugar().Errorf("[%v] [%v] %v", c.ClientIP(), reqmsg, err)
				resp.Msg = Status_500_db
				c.JSON(http.StatusInternalServerError, resp)
				return
			}
			// Send verification code to email
			body := "Hello, " + reqmsg.Mailbox + "!\nWelcome to CESS-GATEWAY authentication service, please write captcha to the authorization page.\ncaptcha: "
			body += fmt.Sprintf("%v", captcha)
			body += "\nValidity: 5 minutes"
			err = communication.SendPlainMail(
				configs.Confile.EmailHost,
				configs.Confile.EmailHostPort,
				configs.Confile.EmailAddress,
				configs.Confile.EmailPassword,
				[]string{reqmsg.Mailbox},
				configs.EmailSubject_captcha,
				body,
			)
			if err != nil {
				Err.Sugar().Errorf("[%v] [%v] %v", c.ClientIP(), reqmsg, err)
				resp.Code = http.StatusBadRequest
				resp.Msg = Status_400_EmailSmpt
				c.JSON(http.StatusBadRequest, resp)
				return
			}
			resp.Code = http.StatusOK
			resp.Msg = Status_200_default
			c.JSON(http.StatusOK, resp)
			return
		}
		Err.Sugar().Errorf("[%v] [%v] %v", c.ClientIP(), reqmsg, err)
		resp.Msg = Status_500_db
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	v := strings.Split(string(bytes), "#")
	if len(v) == 2 {
		vi, err := strconv.ParseInt(v[1], 10, 64)
		if err != nil {
			Err.Sugar().Errorf("[%v] [%v] %v", c.ClientIP(), reqmsg, err)
			resp.Msg = Status_500_unexpected
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
		if time.Now().Unix() >= time.Unix(vi, 0).Unix() {
			Out.Sugar().Infof("[%v] [%v] Captcha has expired and a new captcha has been sent to your mailbox", c.ClientIP(), reqmsg)
			captcha := tools.RandomInRange(100000, 999999)
			v := fmt.Sprintf("%v", captcha) + "#" + fmt.Sprintf("%v", time.Now().Add(time.Minute*10).Unix())
			err = db.Put([]byte(reqmsg.Mailbox), []byte(v))
			if err != nil {
				Err.Sugar().Errorf("[%v] [%v] %v", c.ClientIP(), reqmsg, err)
				resp.Msg = Status_500_db
				c.JSON(http.StatusInternalServerError, resp)
				return
			}
			// Send verification code to email
			body := "Hello, " + reqmsg.Mailbox + "!\nWelcome to CESS-GATEWAY authentication service, please write captcha to the authorization page.\ncaptcha: "
			body += fmt.Sprintf("%v", captcha)
			body += "\nValidity: 5 minutes"
			err = communication.SendPlainMail(
				configs.Confile.EmailHost,
				configs.Confile.EmailHostPort,
				configs.Confile.EmailAddress,
				configs.Confile.EmailPassword,
				[]string{reqmsg.Mailbox},
				configs.EmailSubject_captcha,
				body,
			)
			if err != nil {
				Err.Sugar().Errorf("[%v] [%v] %v", c.ClientIP(), reqmsg, err)
				resp.Code = http.StatusBadRequest
				resp.Msg = Status_400_EmailSmpt
				c.JSON(http.StatusBadRequest, resp)
				return
			}

			resp.Code = http.StatusOK
			resp.Msg = Status_200_expired
			c.JSON(http.StatusOK, resp)
			return
		}
		vi, err = strconv.ParseInt(v[0], 10, 32)
		if err != nil {
			Err.Sugar().Errorf("[%v] [%v] %v", c.ClientIP(), reqmsg, err)
			resp.Msg = Status_500_unexpected
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
		if reqmsg.Captcha != vi {
			Err.Sugar().Errorf("[%v] [%v] %v", c.ClientIP(), reqmsg, err)
			resp.Msg = Status_400_captcha
			c.JSON(http.StatusBadRequest, resp)
			return
		}
		// Send token to user mailbox
		usertoken, err := token.GenerateNewToken(reqmsg.Mailbox)
		if err != nil {
			Err.Sugar().Errorf("[%v] [%v] %v", c.ClientIP(), reqmsg, err)
			resp.Msg = Status_500_unexpected
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
		err = db.Put([]byte(reqmsg.Mailbox), []byte(usertoken))
		if err != nil {
			Err.Sugar().Errorf("[%v] [%v] %v", c.ClientIP(), reqmsg, err)
			resp.Msg = Status_500_db
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
		body := "Hello, " + reqmsg.Mailbox + "!\nCongratulations on your successful authentication, your token is:\n"
		body += usertoken
		fmt.Println(body)
		err = communication.SendPlainMail(
			configs.Confile.EmailHost,
			configs.Confile.EmailHostPort,
			configs.Confile.EmailAddress,
			configs.Confile.EmailPassword,
			[]string{reqmsg.Mailbox},
			configs.EmailSubject_token,
			body,
		)
		if err != nil {
			Err.Sugar().Errorf("[%v] [%v] %v", c.ClientIP(), reqmsg, err)
			resp.Code = http.StatusBadRequest
			resp.Msg = Status_400_EmailSmpt
			c.JSON(http.StatusBadRequest, resp)
			return
		}
		resp.Code = http.StatusOK
		resp.Msg = Status_200_default
		c.JSON(http.StatusOK, resp)
		return
	}

	b, err := token.DecryptToken(string(bytes))
	if err != nil {
		Err.Sugar().Errorf("[%v] [%v] %v", c.ClientIP(), reqmsg, err)
		resp.Msg = Status_500_unexpected
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	var utoken token.TokenMsgType
	err = json.Unmarshal(b, &utoken)
	if err != nil {
		Err.Sugar().Errorf("[%v] [%v] %v", c.ClientIP(), reqmsg, err)
		resp.Msg = Status_500_unexpected
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	if time.Now().Unix() < utoken.ExpirationTime {
		resp.Code = http.StatusOK
		resp.Msg = Status_200_default
		resp.Data = "token=" + string(bytes)
		c.JSON(http.StatusOK, resp)
		return
	}

	newtoken, err := token.RefreshToken(utoken)
	if err != nil {
		Err.Sugar().Errorf("[%v] [%v] %v", c.ClientIP(), reqmsg, err)
		resp.Msg = Status_500_unexpected
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	err = db.Put([]byte(utoken.Mailbox), []byte(newtoken))
	if err != nil {
		Err.Sugar().Errorf("[%v] [%v] %v", c.ClientIP(), reqmsg, err)
		resp.Msg = Status_500_db
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	resp.Code = http.StatusOK
	resp.Msg = Status_200_default
	resp.Data = "token=" + newtoken
	c.JSON(http.StatusOK, resp)
	return
}

// func RegrantTokenHandler(c *gin.Context) {
// 	var resp = RespMsg{
// 		Code: http.StatusBadRequest,
// 		Msg:  "",
// 	}
// 	usertoken_en := c.PostForm("token")
// 	bytes, err := token.DecryptToken(usertoken_en)
// 	if err != nil {
// 		resp.Msg = "illegal token"
// 		c.JSON(http.StatusBadRequest, resp)
// 		return
// 	}
// 	var usertoken token.TokenMsgType
// 	err = json.Unmarshal(bytes, &usertoken)
// 	if err != nil {
// 		resp.Msg = "token format error"
// 		c.JSON(http.StatusBadRequest, resp)
// 		return
// 	}

// 	if time.Now().Add(-(time.Hour * 24 * 3)).Unix() > usertoken.Expire {
// 		resp.Code = http.StatusForbidden
// 		resp.Msg = "The token has expired more than 3 days"
// 		c.JSON(http.StatusForbidden, resp)
// 		return
// 	}

// 	resp.Code = http.StatusInternalServerError
// 	db, err := db.GetDB()
// 	if err != nil {
// 		Err.Sugar().Errorf("[%v] [%v] %v", usertoken.Blocknumber, usertoken.Walletaddr, err)
// 		resp.Msg = err.Error()
// 		c.JSON(http.StatusInternalServerError, resp)
// 		return
// 	}
// 	expire := time.Now().Add(time.Hour * 24 * 7).Unix()
// 	tk, err := token.GetToken(usertoken.Walletaddr, usertoken.Blocknumber, usertoken.Userid, expire)
// 	if err != nil {
// 		Err.Sugar().Errorf("[%v] [%v] %v", usertoken.Blocknumber, usertoken.Walletaddr, err)
// 		resp.Msg = err.Error()
// 		c.JSON(http.StatusInternalServerError, resp)
// 		return
// 	}
// 	//store token to database
// 	err = db.Put([]byte(usertoken.Walletaddr+"_token"), []byte(tk))
// 	if err != nil {
// 		Err.Sugar().Errorf("[%v] [%v] %v", usertoken.Blocknumber, usertoken.Walletaddr, err)
// 		resp.Msg = err.Error()
// 		c.JSON(http.StatusInternalServerError, resp)
// 		return
// 	}
// 	resp.Code = 200
// 	resp.Msg = "success"
// 	resp.Data = tk
// 	c.JSON(http.StatusOK, resp)
// 	return
// }
