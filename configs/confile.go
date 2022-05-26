package configs

import (
	"cess-gateway/tools"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

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
