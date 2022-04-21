package configs

import (
	"errors"
	"os"

	"github.com/spf13/viper"
)

type Configfile struct {
	ChainAddr     string `toml:"ChainAddr"`
	ServiceAddr   string `toml:"ServiceAddr"`
	ServicePort   string `toml:"ServicePort"`
	AccountAddr   string `toml:"AccountAddr"`
	AccountSeed   string `toml:"AccountSeed"`
	EmailAddress  string `toml:"EmailAddress"`
	EmailPassword string `toml:"EmailPassword"`
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
	return nil
}
