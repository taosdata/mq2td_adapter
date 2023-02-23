//go:build !windows
// +build !windows

package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "default",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Init()
			assert.Equal(t, &Config{
				TDengine: TDengine{
					Host:     "127.0.0.1",
					Port:     6030,
					User:     "root",
					Password: "taosdata",
					DB:       "test_mqtt",
				},
				MQTT: MQTT{
					Address:   "tcp://127.0.0.1:1883",
					ClientID:  "",
					Username:  "",
					Password:  "",
					KeepAlive: 30,
					CAPath:    "",
					CertPath:  "",
					KeyPath:   "",
				},
				ShowSql:  false,
				LogLevel: "info",
				RulePath: "./config/rule.json",
				Log: Log{
					Path:          "/var/log/taos",
					RotationCount: 7,
					RotationTime:  time.Hour * 24,
					RotationSize:  1 * 1024 * 1024 * 1024, // 1G
				},
			}, Conf)
		})
	}
}
