package main

import (
	"fmt"

	"github.com/Eminlin/adbCtrlGo/app/adb"
)

func main() {
	client, err := adb.NewAdbClient(adb.AdbClient{
		IP:    "127.0.0.1",
		Port:  5555,
		Debug: true,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(client.GetPhoneModel())
	fmt.Println(client.GetAdbVersion())
	fmt.Println(client.Ping("www.baidu.com"))
}
