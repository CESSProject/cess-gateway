package token

import (
	"cess-httpservice/internal/encryption"
	"encoding/base64"
	"encoding/json"
)

type TokenMsgType struct {
	Walletaddr  string `json:"walletaddr"`
	Blocknumber int64  `json:"blocknumber"`
	Expire      int64  `json:"expire"`
}

// Generate user token
func GetToken(walletaddr string, blocknumber, expire int64) (string, error) {
	token := TokenMsgType{
		Walletaddr:  walletaddr,
		Blocknumber: blocknumber,
		Expire:      expire,
	}
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
