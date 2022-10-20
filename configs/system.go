package configs

import (
	"time"
)

// type and version
const VERSION = "CESS-Gateway v0.2.1.221019.1500"

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

	// Valid Time Of Captcha
	ValidTimeOfCaptcha = time.Duration(time.Minute * 5)

	//
	SIZE_1KB = 1024
	SIZE_1MB = 1024 * SIZE_1KB
	SIZE_1GB = 1024 * SIZE_1MB
)

const (
	//Scheduler's rpc service name
	RpcService_Scheduler = "wservice"
	//Scheduler's rpc service name
	RpcService_Miner = "mservice"
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

const (
	HELP_common = `Please check with the following help information:
    1.Check if the wallet balance is sufficient
    2.Block hash:`
	HELP_BuySpace1 = `Please check with the following help information:
    1.Check whether the available space is sufficient
    2.Check if the wallet balance is sufficient
    3.Block hash:`
	HELP_BuySpace2 = `    4.Check the fileBank.buyPackage transaction event result in the block hash above:
        If system.ExtrinsicFailed is prompted, it means failure;
        If system.ExtrinsicSuccess is prompted, it means success;`
	HELP_Upgrade = `    3.Check the fileBank.upgradePackage transaction event result in the block hash above:
        If system.ExtrinsicFailed is prompted, it means failure;
        If system.ExtrinsicSuccess is prompted, it means success;`
	HELP_Renewal = `    3.Check the fileBank.renewalPackage transaction event result in the block hash above:
        If system.ExtrinsicFailed is prompted, it means failure;
        If system.ExtrinsicSuccess is prompted, it means success;`
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

var PublicKey []byte
