package pool

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBytesPoolGet(t *testing.T) {
	b := BytesPoolGet()
	b.WriteByte('a')
	s := b.String()
	assert.Equal(t, "a", s)
	BytesPoolPut(b)
}
