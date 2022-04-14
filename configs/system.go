package configs

import "time"

const (
	// version
	VERSION = "cess-httpservice_v0.0.0"

	// CESS chain addr
	ChainAddr = "ws://106.15.44.155:9948/"

	// base dir
	BaseDir = "/usr/local/cess"

	// log file dir
	LogfileDir = BaseDir + "/log"

	// keyfile dir
	PrivateKeyfile = BaseDir + "/.privateKey.pem"
	PublicKeyfile  = BaseDir + "/.publicKey.pem"

	// database dir
	DbDir = BaseDir + "/db"

	// file cache dir
	FileCacheDir = BaseDir + "/cache"

	// http service port
	ServicePort = "8081"

	// random number valid time, the unit is minutes
	RandomValidTime = 5.0

	// the time to wait for the event, in seconds
	TimeToWaitEvents = time.Duration(time.Second * 15)

	//The minimum deposit when the user is working normally
	MinimumDeposit = "10000000000000"

	//The minimum deposit when the user is working normally
	CessTokenAccuracy = "1000000000000"
)

const (
	RpcService_Scheduler = "wservice"
	RpcMethod_WriteFile  = "writefile"
	RpcMethod_ReadFile   = "readfile"
)
