package adb

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Eminlin/adbCtrlGo/app/format"
	"github.com/beevik/etree"
)

func ParseDump() {
	doc := etree.NewDocument()
	err := doc.ReadFromFile("temp/dump.xml")
	if err != nil {
		fmt.Println(err)
		return
	}
	root := doc.SelectElement("hierarchy")
	for _, v := range root.FindElements("") {
		attr(v)
	}
}

func attr(v *etree.Element) {
	// fmt.Printf("%+v\n", v)
	for _, k := range v.Attr {
		//头像
		if k.Value == "com.ss.android.ugc.aweme:id/user_avatar" {
			fmt.Println(k.Element().SelectAttr("content-desc").Value)
		}
		//用户名
		if k.Value == "com.ss.android.ugc.aweme:id/title" {
			fmt.Println(k.Element().SelectAttr("text").Value)
		}
		//点赞
		if k.Value == "com.ss.android.ugc.aweme:id/c0m" {
			fmt.Println(k.Element().SelectAttr("content-desc").Value)
			parseBounds(k.Element().SelectAttr("bounds").Value)
		}
		//评论
		if k.Value == "com.ss.android.ugc.aweme:id/b8h" {
			fmt.Println(k.Element().SelectAttr("content-desc").Value)
		}
		//收藏
		if k.Value == "com.ss.android.ugc.aweme:id/b3k" {
			fmt.Println(k.Element().SelectAttr("content-desc").Value)
		}
		//分享
		if k.Value == "com.ss.android.ugc.aweme:id/l5s" {
			fmt.Println(k.Element().SelectAttr("content-desc").Value)
		}
		//歌名
		if k.Value == "com.ss.android.ugc.aweme:id/is7" {
			fmt.Println(k.Element().SelectAttr("content-desc").Value)
		}
		//描述
		if k.Value == "com.ss.android.ugc.aweme:id/desc" {
			fmt.Println(k.Element().SelectAttr("text").Value)
		}
	}
}

//	[36,1410][301,1485]
func parseBounds(b string) format.ButtonPoint {
	rtn := format.ButtonPoint{}
	if b == "" {
		return rtn
	}
	b = strings.ReplaceAll(b, "][", ",")
	b = strings.ReplaceAll(b, "[", "")
	b = strings.ReplaceAll(b, "]", "")
	bTemp := strings.Split(b, ",")
	if len(bTemp) != 4 {
		fmt.Println(bTemp)
		return rtn
	}
	temp := []int{}
	for _, v := range bTemp {
		t, _ := strconv.Atoi(v)
		temp = append(temp, t)
	}
	rtn.XRangeL = temp[0]
	rtn.XRangeR = temp[2]
	rtn.YRangeL = temp[1]
	rtn.YRangeR = temp[3]
	return rtn
}
