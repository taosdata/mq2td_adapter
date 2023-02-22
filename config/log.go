package config

import (
	"runtime"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Log struct {
	Path          string
	RotationCount uint
	RotationTime  time.Duration
	RotationSize  uint
}

func initLog() {
	switch runtime.GOOS {
	case "windows":
		viper.SetDefault("log.path", "C:\\TDengine\\log")
		_ = viper.BindEnv("log.path", "MQ2TD_LOG_PATH")
		pflag.String("log.path", "C:\\TDengine\\log", `log path. Env "MQ2TD_LOG_PATH"`)
	default:
		viper.SetDefault("log.path", "/var/log/taos")
		_ = viper.BindEnv("log.path", "MQ2TD_LOG_PATH")
		pflag.String("log.path", "/var/log/taos", `log path. Env "MQ2TD_LOG_PATH"`)
	}

	viper.SetDefault("log.rotationCount", 7)
	_ = viper.BindEnv("log.rotationCount", "MQ2TD_LOG_ROTATION_COUNT")
	pflag.Uint("log.rotationCount", 7, `log rotation count. Env "MQ2TD_LOG_ROTATION_COUNT"`)

	viper.SetDefault("log.rotationTime", time.Hour*24)
	_ = viper.BindEnv("log.rotationTime", "MQ2TD_LOG_ROTATION_TIME")
	pflag.Duration("log.rotationTime", time.Hour*24, `log rotation time. Env "MQ2TD_LOG_ROTATION_TIME"`)

	viper.SetDefault("log.rotationSize", "1GB")
	_ = viper.BindEnv("log.rotationSize", "MQ2TD_LOG_ROTATION_SIZE")
	pflag.String("log.rotationSize", "1GB", `log rotation size(KB MB GB), must be a positive integer. Env "MQ2TD_LOG_ROTATION_SIZE"`)
}

func (l *Log) setValue() {
	l.Path = viper.GetString("log.path")
	l.RotationCount = viper.GetUint("log.rotationCount")
	l.RotationTime = viper.GetDuration("log.rotationTime")
	l.RotationSize = viper.GetSizeInBytes("log.rotationSize")
}
