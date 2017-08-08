package evmn

import (
	"expvar"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	code200 = expvar.NewInt("code:200")

	expected = map[string]string{
		"list":        "cmdline code memstats",
		"config code": "graph_title code\ngraph_category expvar\ngraph_args --base 1000 --units=si\n_200.label 200\n_200.min 0\n_200.type DERIVE\n.",
		"fetch code":  "_200.value 42\n.",
	}
)

func init() {
	code200.Add(42)
}

func TestHandler(t *testing.T) {
	for k, v := range expected {
		r, err := handler(k)
		assert.NoError(t, err)
		assert.Equal(t, v, r)
	}
}
