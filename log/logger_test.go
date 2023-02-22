package log

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/taosdata/mq2td_adapter/config"
)

func Test_defaultPool_Put(t *testing.T) {
	b := bufferPool.Get()
	b.WriteByte('a')
	s := b.String()
	assert.Equal(t, "a", s)
	bufferPool.Put(b)
}

func TestGetLogger(t *testing.T) {
	config.Init()
	ConfigLog()
	l := GetLogger("test")
	l.WithField("test", "test1").Info("test")
	Close(context.Background())
}
