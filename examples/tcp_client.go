package main

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/samhug/go-libzt"
	"io/ioutil"
	"os"
	"time"
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
	fmt.Printf("My Address [%s]\n", addr)

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
