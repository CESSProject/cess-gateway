package configs

import "time"

const (
	// base dir
	BaseDir = "/usr/local/cess-gateway"

	// log file dir
	LogfileDir = BaseDir + "/log"

	// keyfile dir
	PrivateKeyfile = BaseDir + "/.privateKey.pem"
	PublicKeyfile  = BaseDir + "/.publicKey.pem"

	// database dir
	DbDir = BaseDir + "/db"

	// file cache dir
	FileCacheDir = BaseDir + "/cache"

	// file records dir
	FilRecordsDir = "records"

	// random number valid time, the unit is minutes
	RandomValidTime = 5.0

	// the time to wait for the event, in seconds
	TimeToWaitEvents = time.Duration(time.Second * 20)

	// The validity period of the token, the default is 30 days
	ValidTimeOfToken = time.Duration(time.Hour * 24 * 30)

	//
	SIZE_1GB = 1024 * 1024 * 1024
)

const (
	//Scheduler's rpc service name
	RpcService_Scheduler = "wservice"
	//auth method of rpc service
	RpcMethod_auth = "auth"
	//write method of rpc service
	RpcMethod_WriteFile = "writefile"
	//read method of rpc service
	RpcMethod_ReadFile = "readfile"
	//
	RpcBuffer = 1024 * 1024

	//
	EmailSubject_captcha = "CESS | Authorization captcha"
	EmailSubject_token   = "CESS | Authorization token"
)

// return state code
const (
	Code_200 = 200
	Code_400 = 400
	Code_403 = 403
	Code_404 = 404
	Code_500 = 500
	Code_600 = 600
)
