package config

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type MQTT struct {
	Address   string `json:"address"`
	ClientID  string `json:"clientID"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	KeepAlive int64  `json:"keepAlive"`
	CAPath    string `json:"caPath"`
	CertPath  string `json:"certPath"`
	KeyPath   string `json:"keyPath"`
}

func initMQTT() {
	viper.SetDefault("mqtt.address", "tcp://127.0.0.1:1883")
	_ = viper.BindEnv("mqtt.address", "MQ2TD_MQTT_ADDRESS")
	pflag.String("mqtt.address", "tcp://127.0.0.1:1883", `mqtt address. Env "MQ2TD_MQTT_ADDRESS"`)

	viper.SetDefault("mqtt.clientID", "")
	_ = viper.BindEnv("mqtt.clientID", "MQ2TD_MQTT_CLIENT_ID")
	pflag.String("mqtt.clientID", "", `mqtt clientID. Env "MQ2TD_MQTT_CLIENT_ID"`)

	viper.SetDefault("mqtt.username", "")
	_ = viper.BindEnv("mqtt.username", "MQ2TD_MQTT_USERNAME")
	pflag.String("mqtt.username", "", `mqtt username. Env "MQ2TD_MQTT_USERNAME"`)

	viper.SetDefault("mqtt.password", "")
	_ = viper.BindEnv("mqtt.password", "MQ2TD_MQTT_PASSWORD")
	pflag.String("mqtt.password", "", `mqtt password. Env "MQ2TD_MQTT_PASSWORD"`)

	viper.SetDefault("mqtt.keepAlive", int64(30))
	_ = viper.BindEnv("mqtt.keepAlive", "MQ2TD_MQTT_KEEP_ALIVE")
	pflag.Int64("mqtt.keepAlive", int64(30), `mqtt password. Env "MQ2TD_MQTT_KEEP_ALIVE"`)

	viper.SetDefault("mqtt.caPath", "")
	_ = viper.BindEnv("mqtt.caPath", "MQ2TD_MQTT_CA_PATH")
	pflag.String("mqtt.caPath", "", `mqtt ca file path. Env "MQ2TD_MQTT_CA_PATH"`)

	viper.SetDefault("mqtt.certPath", "")
	_ = viper.BindEnv("mqtt.certPath", "MQ2TD_MQTT_CERT_PATH")
	pflag.String("mqtt.certPath", "", `mqtt cert file path. Env "MQ2TD_MQTT_CERT_PATH"`)

	viper.SetDefault("mqtt.keyPath", "")
	_ = viper.BindEnv("mqtt.keyPath", "MQ2TD_MQTT_KEY_PATH")
	pflag.String("mqtt.keyPath", "", `mqtt key file path. Env "MQ2TD_MQTT_KEY_PATH"`)

}

func (m *MQTT) setValue() {
	m.Address = viper.GetString("mqtt.address")
	m.ClientID = viper.GetString("mqtt.clientID")
	m.Username = viper.GetString("mqtt.username")
	m.Password = viper.GetString("mqtt.password")
	m.KeepAlive = viper.GetInt64("mqtt.keepAlive")
	m.CAPath = viper.GetString("mqtt.caPath")
	m.CertPath = viper.GetString("mqtt.certPath")
	m.KeyPath = viper.GetString("mqtt.keyPath")
}
