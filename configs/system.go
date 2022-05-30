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
	TimeToWaitEvents = time.Duration(time.Second * 15)

	// The validity period of the token, the default is 30 days
	ValidTimeOfToken = time.Duration(time.Hour * 24 * 30)
)

const (
	//Scheduler's rpc service name
	RpcService_Scheduler = "wservice"
	//write method of rpc service
	RpcMethod_WriteFile = "writefile"
	//read method of rpc service
	RpcMethod_ReadFile = "readfile"
	//
	RpcBuffer = 64 * 1024

	//
	EmailSubject_captcha = "CESS | Authorization captcha"
	EmailSubject_token   = "CESS | Authorization token"
)
