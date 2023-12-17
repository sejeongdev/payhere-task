package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// ViperConfig ...
type ViperConfig struct {
	*viper.Viper
}

// Payhere ...
var Payhere *ViperConfig

var (
	ConfServerPORT int    = 18080
	ENVServerLOCAL string = ".env.local"
)

func init() {
	pflag.BoolP("version", "v", false, "Show version number and quit")
	pflag.IntP("port", "p", ConfServerPORT, "server Port")
	pflag.Parse()

	var err error
	Payhere, err = readConfig(getDefaultConfig())
	if err != nil {
		os.Exit(1)
	}

	Payhere.BindPFlags(pflag.CommandLine)
}

func getDefaultConfig() map[string]any {
	return map[string]any{
		"port":      ConfServerPORT,
		"env":       "local",
		"db_prefix": "payhere",
	}
}

func readConfig(defaults map[string]any) (*ViperConfig, error) {
	v := viper.New()
	for key, value := range defaults {
		v.SetDefault(key, value)
	}
	v.AddConfigPath("./")
	v.AddConfigPath("./config")
	v.AddConfigPath("../")
	v.AddConfigPath("../config")
	v.AddConfigPath("../../config")
	v.AutomaticEnv()

	switch strings.ToLower(v.GetString("ENV")) {
	case "local":
		v.SetConfigName(ENVServerLOCAL)
		v.Debug()
	}
	err := v.ReadInConfig()
	switch err.(type) {
	default:
		fmt.Println("error ", err)
		return &ViperConfig{}, err
	case nil:
		break
	case viper.ConfigFileNotFoundError:
		fmt.Printf("Warn: %s\n", err)
	}
	return &ViperConfig{
		Viper: v,
	}, nil
}
