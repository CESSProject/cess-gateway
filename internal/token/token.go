package token

import (
	"cess-httpservice/configs"
	"cess-httpservice/internal/encryption"
	"cess-httpservice/tools"
	"encoding/base64"
	"encoding/json"
	"time"
)

// type TokenMsgType struct {
// 	Userid      int64  `json:"userid"`
// 	Blocknumber uint64 `json:"blocknumber"`
// 	Expire      int64  `json:"expire"`
// 	Walletaddr  string `json:"walletaddr"`
// 	Randomcode  string `json:"randomcode"`
// }

type TokenMsgType struct {
	UserId          int64  `json:"userId"`
	CreateUserTime  int64  `json:"createUserTime"`
	CreateTokenTime int64  `json:"createTokenTime"`
	ExpirationTime  int64  `json:"expirationTime"`
	Mailbox         string `json:"mailbox"`
	RandomCode      string `json:"randomCode"`
}

// Generate a new token
func GenerateNewToken(mailbox string) (string, error) {
	var (
		err   error
		token = TokenMsgType{}
	)
	token.UserId, err = tools.GetGuid(int64(tools.RandomInRange(0, 1023)))
	if err != nil {
		return "", err
	}
	token.RandomCode = tools.GetRandomcode(16)
	token.Mailbox = mailbox
	t := time.Now().Unix()
	token.CreateUserTime = t
	token.CreateTokenTime = t
	token.ExpirationTime = time.Unix(t, 0).Add(configs.ValidTimeOfToken).Unix()
	bytes, err := json.Marshal(token)
	if err != nil {
		return "", err
	}

	en, err := encryption.RSA_Encrypt(bytes)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(en), nil
}

// Decode user token
func DecryptToken(token string) ([]byte, error) {
	en, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return nil, err
	}
	bytes, err := encryption.RSA_Decrypt(en)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
