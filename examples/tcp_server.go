package main

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/samhug/go-libzt"
	"io/ioutil"
	"os"
)

func main() {
	const EARTH = "8056c2e21c000001"
	const PORT uint16 = 8888

	// Create a temporary directory to use as ZeroTier's home dir
	ztHomePath, err := ioutil.TempDir("", "zt-home-")
	if err != nil {
		glog.Fatal(err)
	}
	defer os.RemoveAll(ztHomePath)

	fmt.Printf("Using %s as temporary ZeroTier home path\n", ztHomePath)

	zt := libzt.Init(EARTH, ztHomePath)

	addr, err := zt.GetIPv6Address()
	if err != nil {
		glog.Fatal(err)
	}
	fmt.Printf("Listening at [%s]:%d\n", addr, PORT)

	conn, err := zt.Listen6(PORT)
	if err != nil {
		glog.Fatal(err)
	}

	fmt.Println("Waiting for connection")

	tcpConn, err := conn.Accept()
	if err != nil {
		glog.Fatal(err)
	}

	fmt.Println("Accepted connection")

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
