package evmn

import (
	"fmt"
	"math/rand"
	"net"
	"net/textproto"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	hostname = fmt.Sprint(rand.Int63())
	expected["nodes"] = hostname + "\n."
}

func TestListenAndServe(t *testing.T) {
	go func() {
		err := ListenAndServe(":4950")
		assert.NoError(t, err)
	}()
}

func TestServer(t *testing.T) {
	expected["config"] = "# Unknown service"
	expected["fetch"] = "# Unknown service"
	expected["help"] = "# Unknown command"

	for k, v := range expected {
		conn, err := net.Dial("tcp", "localhost:4950")
		assert.NoError(t, err)

		tc := textproto.NewConn(conn)
		s, err := tc.ReadLine()
		assert.NoError(t, err)
		assert.Equal(t, "# munin node at "+hostname, s)

		n, err := tc.W.WriteString(k + "\n\n")
		assert.NoError(t, tc.W.Flush())
		assert.NoError(t, err)
		assert.NotZero(t, n)
		r, err := tc.ReadLine()
		assert.NoError(t, err)
		assert.Equal(t, strings.Split(v, "\n")[0], r)
		conn.Close()
	}
}
