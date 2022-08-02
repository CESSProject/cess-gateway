package configs

import (
	"cess-gateway/tools"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// type and version
const VERSION = "CESS-Gateway v0.1.2.220802.1632"

type Configfile struct {
	RpcAddr       string `toml:"RpcAddr"`
	ServiceAddr   string `toml:"ServiceAddr"`
	ServicePort   string `toml:"ServicePort"`
	AccountSeed   string `toml:"AccountSeed"`
	EmailAddress  string `toml:"EmailAddress"`
	EmailPassword string `toml:"EmailPassword"`
	EmailHost     string `toml:"EmailHost"`
	EmailHostPort int    `toml:"EmailHostPort"`
}

var Confile = new(Configfile)
var ConfigFilePath string

func ParseConfile() error {
	f, err := os.Stat("conf.toml")
	if err != nil {
		return err
	}
	if f.IsDir() {
		return errors.New("conf.toml not found")
	}

	viper.SetConfigFile("conf.toml")
	viper.SetConfigType("toml")

	err = viper.ReadInConfig()
	if err != nil {
		return err
	}

	err = viper.Unmarshal(Confile)
	if err != nil {
		return err
	}

	if !tools.VerifyMailboxFormat(Confile.EmailAddress) {
		fmt.Printf("\x1b[%dm[err]\x1b[0m '%v' email format error\n", 41, Confile.EmailAddress)
		os.Exit(1)
	}

	return nil
}

const ConfigFile_Templete = `#The rpc address of the chain node
RpcAddr           = ""
#The ip address that the cess-gateway service listens to
ServiceAddr       = ""
#The port number on which the cess-gateway service listens
ServicePort       = ""
#Phrase or seed of for wallet account
AccountSeed       = ""
#Email address
EmailAddress      = ""
#Email authorization code
AuthorizationCode = ""
#Outgoing server address of SMTP service
EmailHost         = ""
#Outgoing server port number of SMTP service
EmailHostPort     = 0`
