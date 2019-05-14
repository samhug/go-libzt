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
	udp, err := zt.Listen6UDP(PORT)
	if err != nil {
		glog.Fatal(err)
	}

	fmt.Println("Waiting for messages")

	buffer := make([]byte, 1024)

	for {
		len, addr, err := udp.ReadFrom(buffer)
		fmt.Println("Received: ", string(buffer), len, addr, err)
	}
}
