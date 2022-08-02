package cmd

import (
	"cess-gateway/configs"
	"cess-gateway/internal/chain"
	"cess-gateway/internal/handler"
	"cess-gateway/internal/logger"
	"cess-gateway/tools"
	"fmt"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"strconv"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	Name        = "cess-gateway"
	Description = "Rest API service implementation for accessing CESS."
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   Name,
	Short: Description,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// init
func init() {
	rootCmd.AddCommand(
		Command_Default(),
		Command_Version(),
		Command_Run(),
		Command_BuySpace(),
		Command_UpgradePackage(),
		Command_Renewal(),
	)
	rootCmd.PersistentFlags().StringVarP(&configs.ConfigFilePath, "config", "c", "", "Custom profile")
}

func Command_Version() *cobra.Command {
	cc := &cobra.Command{
		Use:                   "version",
		Short:                 "Print version information",
		Run:                   Command_Version_Runfunc,
		DisableFlagsInUseLine: true,
	}
	return cc
}

func Command_Default() *cobra.Command {
	cc := &cobra.Command{
		Use:                   "default",
		Short:                 "Generate profile template",
		Run:                   Command_Default_Runfunc,
		DisableFlagsInUseLine: true,
	}
	return cc
}

func Command_Run() *cobra.Command {
	cc := &cobra.Command{
		Use:                   "run",
		Short:                 "Operation scheduling service",
		Run:                   Command_Run_Runfunc,
		DisableFlagsInUseLine: true,
	}
	return cc
}

func Command_BuySpace() *cobra.Command {
	cc := &cobra.Command{
		Use:                   "buy",
		Short:                 "Buy space packages:[1, 2, 3, 4, 5]",
		Run:                   Command_BuySpace_Runfunc,
		DisableFlagsInUseLine: true,
	}
	return cc
}

func Command_UpgradePackage() *cobra.Command {
	cc := &cobra.Command{
		Use:                   "upgrade",
		Short:                 "Upgrade a small package to a large package",
		Run:                   Command_UpgradePackage_Runfunc,
		DisableFlagsInUseLine: true,
	}
	return cc
}

func Command_Renewal() *cobra.Command {
	cc := &cobra.Command{
		Use:                   "renewal",
		Short:                 "One-month lease term for additional space package",
		Run:                   Command_Renewal_Runfunc,
		DisableFlagsInUseLine: true,
	}
	return cc
}

// Print version number and exit
func Command_Version_Runfunc(cmd *cobra.Command, args []string) {
	fmt.Println(configs.VERSION)
	os.Exit(0)
}

// Generate configuration file template
func Command_Default_Runfunc(cmd *cobra.Command, args []string) {
	tools.WriteStringtoFile(configs.ConfigFile_Templete, "conf_template.toml")
	pwd, err := os.Getwd()
	if err != nil {
		log.Printf("\x1b[%dm[err]\x1b[0m %v\n", 41, err)
		os.Exit(1)
	}
	path := filepath.Join(pwd, "conf_template.toml")
	log.Printf("[ok] %v\n", path)
	os.Exit(0)
}

// start service
func Command_Run_Runfunc(cmd *cobra.Command, args []string) {
	refreshProfile(cmd)
	logger.Log_Init()
	handler.Main()
}

// buy space package
func Command_BuySpace_Runfunc(cmd *cobra.Command, args []string) {
	if len(os.Args) < 3 {
		log.Println("[err] Please enter the correct package type: [1,2,3,4,5]")
		os.Exit(1)
	}
	count := types.NewU128(*big.NewInt(0))
	p_type, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Println("[err] Please enter the correct package type: [1,2,3,4,5]")
		os.Exit(1)
	}
	if p_type < 1 || p_type > 5 {
		log.Println("[err] Please enter the correct package type: [1,2,3,4,5]")
		os.Exit(1)
	}
	if p_type == 5 {
		if len(os.Args) < 4 {
			log.Println("[err] Please enter the purchased space size (unit: TB)")
			os.Exit(1)
		}
		si, err := strconv.ParseUint(os.Args[3], 10, 64)
		if err != nil {
			log.Println("[err] Please enter a number greater than 5")
			os.Exit(1)
		}
		if si < 5 {
			log.Println("[err] Please enter a number greater than 5")
			os.Exit(1)
		}
		count.SetUint64(si)
	}
	refreshProfile(cmd)
	logger.Log_Init()
	txhash, err := chain.BuySpacePackage(types.U8(p_type), count)
	if txhash == "" {
		log.Printf("[err] Failed purchase: %v\n", err)
		os.Exit(1)
	}
	logger.Out.Sugar().Infof("Space purchased successfully: %v", txhash)
	log.Printf("[ok] success\n")
	os.Exit(0)
}

// buy space package
func Command_UpgradePackage_Runfunc(cmd *cobra.Command, args []string) {
	if len(os.Args) < 3 {
		log.Println("[err] Please enter the correct package type: [1,2,3,4,5]")
		os.Exit(1)
	}
	count := types.NewU128(*big.NewInt(0))
	p_type, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Println("[err] Please enter the correct package type: [1,2,3,4,5]")
		os.Exit(1)
	}
	if p_type < 1 || p_type > 5 {
		log.Println("[err] Please enter the correct package type: [1,2,3,4,5]")
		os.Exit(1)
	}
	if p_type == 5 {
		if len(os.Args) < 4 {
			log.Println("[err] Please enter the purchased space size (unit: TB)")
			os.Exit(1)
		}
		si, err := strconv.ParseUint(os.Args[3], 10, 64)
		if err != nil {
			log.Println("[err] Please enter a number greater than 5")
			os.Exit(1)
		}
		if si < 5 {
			log.Println("[err] Please enter a number greater than 5")
			os.Exit(1)
		}
		count.SetUint64(si)
	}
	refreshProfile(cmd)
	logger.Log_Init()
	txhash, err := chain.UpgradeSpacePackage(types.U8(p_type), count)
	if txhash == "" {
		log.Printf("[err] Upgrade package failed: %v\n", err)
		os.Exit(1)
	}
	logger.Out.Sugar().Infof("Upgrade package successfully: %v", txhash)
	log.Printf("[ok] success\n")
	os.Exit(0)
}

// Increase space package lease term
func Command_Renewal_Runfunc(cmd *cobra.Command, args []string) {
	refreshProfile(cmd)
	logger.Log_Init()
	txhash, err := chain.Renewal()
	if txhash == "" {
		log.Printf("[err] Renewal package failed: %v\n", err)
		os.Exit(1)
	}
	logger.Out.Sugar().Infof("Renewal package successfully: %v", txhash)
	log.Printf("[ok] success\n")
	os.Exit(0)
}

func refreshProfile(cmd *cobra.Command) {
	configpath1, _ := cmd.Flags().GetString("config")
	configpath2, _ := cmd.Flags().GetString("c")
	if configpath1 != "" {
		configs.ConfigFilePath = configpath1
	} else {
		configs.ConfigFilePath = configpath2
	}
	parseProfile()
}

func parseProfile() {
	var (
		err          error
		confFilePath string
	)
	if configs.ConfigFilePath == "" {
		confFilePath = "./conf.toml"
	} else {
		confFilePath = configs.ConfigFilePath
	}
	f, err := os.Stat(confFilePath)
	if err != nil {
		log.Printf("[err] The '%v' file does not exist.\n", confFilePath)
		os.Exit(1)
	}
	if f.IsDir() {
		log.Printf("[err] The '%v' is not a file.\n", confFilePath)
		os.Exit(1)
	}

	viper.SetConfigFile(confFilePath)
	viper.SetConfigType("toml")

	err = viper.ReadInConfig()
	if err != nil {
		log.Printf("[err] The '%v' file type error.\n", 41, confFilePath)
		os.Exit(1)
	}

	err = viper.Unmarshal(configs.C)
	if err != nil {
		log.Printf("[err] Configuration file error, please use the default command to generate a template.\n", 41, confFilePath)
		os.Exit(1)
	}

	if configs.C.RpcAddr == "" ||
		configs.C.AccountSeed == "" ||
		configs.C.EmailAddress == "" ||
		configs.C.AuthorizationCode == "" ||
		configs.C.SMTPHost == "" {
		log.Printf("[err] The configuration file cannot have empty entries.\n")
		os.Exit(1)
	}

	if !tools.VerifyMailboxFormat(configs.C.EmailAddress) {
		fmt.Printf("[err] '%v' email format error\n", 41, configs.C.EmailAddress)
		os.Exit(1)
	}

	port, err := strconv.Atoi(configs.C.ServicePort)
	if err != nil {
		log.Printf("[err] Please fill in the correct 'ServicePort'.\n")
		os.Exit(1)
	}
	if port < 1024 {
		log.Printf("[err] Prohibit the use of system reserved port: %v.\n", port)
		os.Exit(1)
	}
	if port > 65535 {
		log.Printf("[err] The 'ServicePort' cannot exceed 65535.\n")
		os.Exit(1)
	}

	//
	if configs.C.SMTPPort <= 0 {
		log.Printf("[err] The 'SMTPPort' is invalid.\n")
		os.Exit(1)
	}

	//
	if err := tools.CreatDirIfNotExist(configs.BaseDir); err != nil {
		log.Printf("[err] %v\n", err)
		os.Exit(1)
	}

	if err := tools.CreatDirIfNotExist(configs.LogfileDir); err != nil {
		log.Printf("[err] %v\n", err)
		os.Exit(1)
	}

	if err := tools.CreatDirIfNotExist(configs.DbDir); err != nil {
		log.Printf("[err] %v\n", err)
		os.Exit(1)
	}

	if err := tools.CreatDirIfNotExist(configs.FileCacheDir); err != nil {
		log.Printf("[err] %v\n", err)
		os.Exit(1)
	}
}
