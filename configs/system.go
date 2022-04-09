package configs

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

	// http service port
	ServicePort = "8081"

	// random number valid time
	RandomValidTime = 10.0
)
