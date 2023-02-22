package pool

import (
	"sync"

	"github.com/panjf2000/ants/v2"
	"github.com/sirupsen/logrus"
)

var GoPool *ants.Pool
var once sync.Once

func InitGoPool(logger *logrus.Entry, poolSize int) {
	once.Do(func() {
		var err error
		GoPool, err = ants.NewPool(poolSize)
		if err != nil {
			logger.WithError(err).Panic("ants.NewPool")
		}
	})
}
