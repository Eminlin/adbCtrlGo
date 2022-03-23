package main

import (
	"fmt"

	"github.com/Eminlin/adbCtrlGo/app/adb"
)

func main() {
	_, err := adb.NewAdbClient("127.0.0.1", 5555)
	if err != nil {
		fmt.Println(err)
	}
}
