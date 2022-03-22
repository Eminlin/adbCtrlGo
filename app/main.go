package app

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Eminlin/adbCtrlGo/app/adb"
	"github.com/Eminlin/adbCtrlGo/app/format"
	"github.com/Eminlin/adbCtrlGo/app/log"
	"github.com/beevik/etree"
)

var l *log.Log = log.NewLog()
var client *adb.AdbClient

func Run() {
	client = adb.NewAdbClient("192.168.31.43", 5555)
	// path, err := client.GetAppPath(douyin)
	// if err != nil {
	// 	log.NewLog().Errorln(err)
	// 	return
	// }
	// client.RunApp(path)
	// l.Infoln(client.GetAppInfo(douyinPackname))
	for {
		time.Sleep(time.Duration(RandInt64(10, 90)) * time.Second)

		client.Swipe(500, 1200, 520, 320)
		time.Sleep(5 * time.Second)

		err := client.GetElement("dump")
		if err != nil {
			l.Errorln(err)
			continue
		}
		err = client.Down("dump")
		if err != nil {
			l.Errorln(err)
			continue
		}
		ParseDump()
	}
}
func ParseDump() {
	doc := etree.NewDocument()
	err := doc.ReadFromFile("temp/dump.xml")
	if err != nil {
		l.Errorln(err)
		return
	}
	root := doc.SelectElement("hierarchy")
	info := format.VideoInfo{}
	for _, v := range root.FindElements("") {
		attr(v, &info)
	}
	l.Infof("获取到新的视频：%+v", info)
}

func attr(v *etree.Element, i *format.VideoInfo) {
	// fmt.Printf("%+v\n", v)
	for _, k := range v.Attr {
		//检测到更新 以后再说
		if k.Value == "com.ss.android.ugc.aweme:id/g2-" {
			l.Infoln("检测到升级弹窗")
			tabBounds(adb.ParseBounds(k.Element().SelectAttr("bounds").Value))
			continue
		}
		if k.Value == "com.ss.android.ugc.aweme:id/avv" {
			l.Infoln("检测到青少年模式")
			tabBounds(adb.ParseBounds(k.Element().SelectAttr("bounds").Value))
			continue
		}
		//用户名
		if k.Value == "com.ss.android.ugc.aweme:id/title" {
			temp := k.Element().SelectAttr("text").Value
			if temp != "" {
				i.Username = temp
			}
		}
		//点赞
		if k.Value == "com.ss.android.ugc.aweme:id/c0m" {
			temp := k.Element().SelectAttr("content-desc").Value
			if temp != "" {
				i.Like = temp
			}
			// adb.ParseBounds(k.Element().SelectAttr("bounds").Value)
		}
		//评论
		if k.Value == "com.ss.android.ugc.aweme:id/b8h" {
			temp := k.Element().SelectAttr("content-desc").Value
			if temp != "" {
				i.Comment = temp
			}
		}
		//收藏
		if k.Value == "com.ss.android.ugc.aweme:id/b3k" {
			temp := k.Element().SelectAttr("content-desc").Value
			if temp != "" {
				i.Collect = temp
			}
		}
		//分享
		if k.Value == "com.ss.android.ugc.aweme:id/l5s" {
			temp := k.Element().SelectAttr("content-desc").Value
			if temp != "" {
				i.Share = temp
			}
		}
		//歌名
		if k.Value == "com.ss.android.ugc.aweme:id/is7" {
			temp := k.Element().SelectAttr("content-desc").Value
			if temp != "" {
				i.Song = temp
			}
		}
		//描述
		if k.Value == "com.ss.android.ugc.aweme:id/desc" {
			temp := k.Element().SelectAttr("text").Value
			if temp != "" {
				i.Desc = temp
			}
		}
	}
}

func tabBounds(e format.ButtonPoint) {
	x := RandInt64(int64(e.XRangeL), int64(e.XRangeR))
	y := RandInt64(int64(e.YRangeL), int64(e.YRangeR))
	err := client.Click(fmt.Sprintf("%d", x), fmt.Sprintf("%d", y))
	if err != nil {
		l.Errorln(err)
	}
}
func RandInt64(min, max int64) int64 {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return rand.Int63n(max-min) + min
}
