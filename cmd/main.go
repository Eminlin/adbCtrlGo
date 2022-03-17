package main

import (
	"github.com/Eminlin/adbCtrlGo/app"
)

func main() {
	go app.Run()
	select {}
}
