package app

import (
	"github.com/Eminlin/adbCtrlGo/app/adb"
)

func Run() {
	client := adb.NewAdbClient("192.168.31.157", 5555)
	// path, err := client.GetAppPath(douyin)
	// fmt.Println(path)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// client.RunApp(path)
	client.GetElement()
	client.Down(uiautomatorXML)
	// adb.ParseDump()
	// for {
	// 	client.Swipe(500, 1200, 520, 320)
	// 	time.Sleep(5 * time.Second)
	// }
	// for {
	// 	fmt.Println("beats")
	// 	time.Sleep(time.Second * 30)
	// }

}
