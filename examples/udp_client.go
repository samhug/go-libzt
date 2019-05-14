package main

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/samhug/go-libzt"
	"os"
)

func main() {
	const EARTH = "8056c2e21c000001"
	const PORT uint16 = 8888

	zt := libzt.Init(EARTH, "/tmp/test-client")
	fmt.Println("My address: ", zt.GetIPv6Address())
	connection, err := zt.Dial6UDP(os.Args[1], PORT)
	if err != nil {
		glog.Fatal(err)
	}

	len, err := connection.Write([]byte("hello world"))
	fmt.Println("Sent: ", len, err)
	connection.Close()
}
