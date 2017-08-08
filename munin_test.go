package evmn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKK(t *testing.T) {
	assert.Equal(t, "", kk(""))
	assert.Equal(t, "a-b_c_d", kk("a-b_c.d"))
}
