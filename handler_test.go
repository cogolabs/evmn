package evmn

import (
	"expvar"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	anyInt  = expvar.NewInt("any")
	code200 = expvar.NewInt("code:200")
	code302 = expvar.NewInt("code:302")
	blank   = expvar.NewInt("")

	intMap = expvar.NewMap("intMap")

	badType = expvar.NewString("badType")

	expected = map[string]string{
		"cap": "multigraph",

		"list":           " any code intMap",
		"config any":     "graph_title any\ngraph_category expvar\ngraph_args --base 1000 --units=si\nany.label any\nany.min 0\nany.type DERIVE\n.",
		"config code":    "graph_title code\ngraph_category expvar\ngraph_args --base 1000 --units=si\n_200.label 200\n_200.min 0\n_200.type DERIVE\n_302.label 302\n_302.min 0\n_302.type DERIVE\n.",
		"config intMap":  "graph_title intMap\ngraph_category expvar\ngraph_args --base 1000 --units=si\nbig.label big\nbig.min 0\nbig.type DERIVE\nsmall.label small\nsmall.min 0\nsmall.type DERIVE\n.",
		"config badType": "graph_title badType\ngraph_category expvar\ngraph_args --base 1000 --units=si\n.",
		"fetch any":      "any.value 123\n.",
		"fetch code":     "_200.value 42\n_302.value 16\n.",
		"fetch intMap":   "big.value 2017\nsmall.value 101\n.",
		"fetch badType":  "\n.",
	}
)

func init() {
	anyInt.Add(123)
	code200.Add(42)
	code302.Add(16)
	intMap.Add("big", 2017)
	intMap.Add("small", 101)
}

func TestHandler(t *testing.T) {
	for k, v := range expected {
		r, err := handler(k)
		assert.NoError(t, err)
		assert.Equal(t, v, r)
	}

	r, err := handler("")
	assert.Equal(t, ErrUnknownCmd, err)
	assert.Equal(t, "", r)
}
