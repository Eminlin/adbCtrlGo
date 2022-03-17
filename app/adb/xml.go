package adb

import (
	"fmt"

	"github.com/beevik/etree"
)

func ParseDump() {
	doc := etree.NewDocument()
	err := doc.ReadFromFile("temp/dump.xml")
	if err != nil {
		fmt.Println(err)
		return
	}
	// root := doc.SelectElement("hierarchy")
	fmt.Println(doc.ReadFromString("468"))
}
