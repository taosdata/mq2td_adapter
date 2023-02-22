package config

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	TDengine TDengine `json:"TDengine"`
	MQTT     MQTT     `json:"mqtt"`
	ShowSql  bool     `json:"showSql"`
	LogLevel string
	RulePath string `json:"rulePath"`
	Log      Log
}

var (
	Conf *Config
)

func Init() {
	viper.SetConfigType("json")
	viper.SetConfigName("config")
	var configPath *string
	viper.AddConfigPath("./config")
	configPath = pflag.StringP("config", "c", "", "config file")
	help := pflag.Bool("help", false, "Print this help message and exit")
	pflag.Parse()
	if *help {
		fmt.Fprintf(os.Stderr, "Usage of mq2td)adapter :\n")
		pflag.PrintDefaults()
		os.Exit(0)
	}
	if *configPath != "" {
		viper.SetConfigFile(*configPath)
	}
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		panic(err)
	}
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("config file not found")
		} else {
			panic(err)
		}
	}
	Conf = &Config{
		ShowSql:  viper.GetBool("showSql"),
		LogLevel: viper.GetString("logLevel"),
		RulePath: viper.GetString("rc"),
	}
	Conf.Log.setValue()
	Conf.MQTT.setValue()
	Conf.TDengine.setValue()
	return
}

func init() {
	viper.SetDefault("showSql", false)
	_ = viper.BindEnv("showSql", "MQ2TD_SHOW_SQL")
	pflag.Bool("showSql", false, `weather to show sql. Env "MQ2TD_SHOW_SQL"`)

	viper.SetDefault("logLevel", "info")
	_ = viper.BindEnv("logLevel", "MQ2TD_LOG_LEVEL")
	pflag.String("logLevel", "info", `log level (panic fatal error warn warning info debug trace). Env "MQ2TD_LOG_LEVEL"`)

	viper.SetDefault("rc", "./config/rule.json")
	_ = viper.BindEnv("rc", "MQ2TD_RULE_PATH")
	pflag.String("rc", "./config/rule.json", `rule config file. ENV "MQ2TD_RULE_PATH""`)

	initLog()
	initMQTT()
	initTDengine()
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		panic(err)
	}
}
