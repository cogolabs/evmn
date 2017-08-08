package evmn

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/textproto"
	"os"
)

var hostname = ""

func init() {
	hostname, _ = os.Hostname()
	if hostname == "" {
		hostname = "localhost"
	}
}

// ListenAndServe opens a munin TCP port, usually ":4949"
func ListenAndServe(address string) error {
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go serve(conn)
	}
}

func serve(conn net.Conn) {
	tc := textproto.NewConn(conn)
	tc.W.Write([]byte(fmt.Sprintf("# munin node at %s\n", hostname)))
	tc.W.Flush()
	for {
		line, err := tc.ReadLine()
		if err != nil {
			if err != io.EOF {
				log.Println(conn.RemoteAddr(), err)
			}
			break
		}
		response, err := handler(line)
		if err != nil {
			tc.W.Write([]byte("# " + err.Error()))
		} else {
			tc.W.Write([]byte(response))
		}
		tc.W.Write([]byte("\n"))
		tc.W.Flush()
	}
	tc.Close()
}
