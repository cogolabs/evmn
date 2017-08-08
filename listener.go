package evmn

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/textproto"
	"os"
	"time"
)

var hostname = ""

func init() {
	hostname, _ = os.Hostname()
	if hostname == "" {
		hostname = "localhost"
	}
}

// Listener opens a munin TCP port, usually ":4949"
func Listener(address string) {
	var (
		ln  net.Listener
		err error
	)
	for ln == nil {
		ln, err = net.Listen("tcp", address)
		if err != nil {
			log.Println(err)
			time.Sleep(time.Minute)
		}
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go func(tc *textproto.Conn) {
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
		}(textproto.NewConn(conn))
	}
}
