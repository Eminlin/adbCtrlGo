package app

import (
	"time"

	"github.com/Eminlin/adbCtrlGo/app/adb"
)

func Run() {
	client := adb.NewAdbClient("127.0.0.1", 5555)
	for {
		client.Swipe(500, 500, 200, 200)
		time.Sleep(5 * time.Second)
	}
	// for {
	// 	fmt.Println("beats")
	// 	time.Sleep(time.Second * 30)
	// }

}
