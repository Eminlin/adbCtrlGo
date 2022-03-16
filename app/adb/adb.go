package adb

import (
	"os"
	"os/exec"
	"runtime"
)

//getAdbCli return adb cli
func getAdbCli() string {
	if runtime.GOOS == "windows" {
		return "./adb/adb.exe"
	}
	return "adb"
}

//clickEventCode key code event
func eventCode(code string) error {
	cmd := exec.Command(getAdbCli(), "shell", "input", "keyevent", code)
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

//click click screen
func click(x, y string) error {
	cmd := exec.Command(getAdbCli(), "shell", "input", "tap", x, y)
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
