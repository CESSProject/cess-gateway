package token

import (
	"cess-httpservice/internal/encryption"
	"cess-httpservice/tools"
	"encoding/base64"
	"encoding/json"
	"time"
)

type TokenMsgType struct {
	Userid      int64  `json:"userid"`
	Blocknumber uint64 `json:"blocknumber"`
	Expire      int64  `json:"expire"`
	Walletaddr  string `json:"walletaddr"`
	Randomcode  string `json:"randomcode"`
}

// Generate user token
func GetToken(walletaddr string, blocknumber uint64, expire int64) (string, error) {
	token := TokenMsgType{
		Userid:      0,
		Blocknumber: blocknumber,
		Expire:      expire,
		Walletaddr:  walletaddr,
		Randomcode:  "",
	}
	uid, err := tools.GetGuid(time.Now().UnixNano())
	if err != nil {
		return "", err
	}

	token.Userid = uid
	token.Randomcode = tools.GetRandomcode(16)

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