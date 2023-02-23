package config

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type TDengine struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DB       string `json:"db"`
}

func initTDengine() {
	viper.SetDefault("tdengine.host", "127.0.0.1")
	_ = viper.BindEnv("tdengine.host", "MQ2TD_TDENGINE_HOST")
	pflag.String("tdengine.host", "127.0.0.1", `TDengine Host. Env "MQ2TD_TDENGINE_HOST"`)

	viper.SetDefault("tdengine.port", int(6030))
	_ = viper.BindEnv("tdengine.port", "MQ2TD_TDENGINE_PORT")
	pflag.Int("tdengine.port", 6030, `TDengine Port. Env "MQ2TD_TDENGINE_PORT"`)

	viper.SetDefault("tdengine.user", "root")
	_ = viper.BindEnv("tdengine.user", "MQ2TD_TDENGINE_USER")
	pflag.String("tdengine.user", "root", `TDengine User. Env "MQ2TD_TDENGINE_USER"`)

	viper.SetDefault("tdengine.password", "taosdata")
	_ = viper.BindEnv("tdengine.password", "MQ2TD_TDENGINE_PASSWORD")
	pflag.String("tdengine.password", "taosdata", `TDengine Password. Env "MQ2TD_TDENGINE_PASSWORD"`)

	viper.SetDefault("tdengine.db", "test_mqtt")
	_ = viper.BindEnv("tdengine.db", "MQ2TD_TDENGINE_DB")
	pflag.String("tdengine.db", "test_mqtt", `TDengine DB. Env "MQ2TD_TDENGINE_DB"`)
}

func (t *TDengine) setValue() {
	t.Host = viper.GetString("tdengine.host")
	t.Port = viper.GetInt("tdengine.port")
	t.User = viper.GetString("tdengine.user")
	t.Password = viper.GetString("tdengine.password")
	t.DB = viper.GetString("tdengine.db")
}
