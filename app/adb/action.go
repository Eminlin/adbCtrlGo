package adb

//ClickHome click home menu
func ClickHome() error {
	return eventCode("3")
}

//ClickBack click back menu
func ClickBack() error {
	return eventCode("4")
}
