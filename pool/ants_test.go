package pool

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestInitGoPool(t *testing.T) {
	InitGoPool(logrus.New().WithField("test", "ants"), 50000)
}
