package main

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/samhug/go-libzt"
)

func main() {
	const EARTH = "8056c2e21c000001"
	const PORT uint16 = 8888

	zt := libzt.Init(EARTH, "/tmp/test-server")
	fmt.Println("My address: ", zt.GetIPv6Address())
	conn, err := zt.Listen6(PORT)
	if err != nil {
		glog.Fatal(err)
	}

	fmt.Println("Waiting for conn")

	tcpConn, err := conn.Accept()
	if err != nil {
		glog.Fatal(err)
	}

	fmt.Println("Accepted conn")

	buffer := make([]byte, 1024)

	for {
		len, err := tcpConn.Read(buffer)
		fmt.Println("Received: ", string(buffer), len, err)
		if err != nil {
			glog.Fatal(err)
		}
		if len == 0 {
			break
		}
	}
}
