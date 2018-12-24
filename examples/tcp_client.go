package main

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/samuelhug/go-libzt"
	"os"
	"time"
)

func main() {
	const EARTH = "8056c2e21c000001"
	const PORT uint16 = 8888

	zt := libzt.Init(EARTH, "/tmp/test-client")
	fmt.Println("My address: ", zt.GetIPv6Address())
	connection, err := zt.Connect6(os.Args[1], PORT)
	fmt.Println("Connected")

	if err != nil {
		glog.Fatal(err)
	}

	len, err := connection.Write([]byte("hello world\n"))
	fmt.Println("Sent: ", len, err)

	time.Sleep(2 * time.Second)
	connection.Close()
}
