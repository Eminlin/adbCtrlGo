package adb

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

const AdbServerPort = 5037
const AdbDaemonPort = 5555

//getAdbCli return adb cli
func getAdbCli() string {
	if runtime.GOOS == "windows" {
		return "adb/win/adb.exe"
	}
	return "adb"
}

//cmdRun exec enter
func cmdRun(cmd *exec.Cmd) (err error) {
	fmt.Printf("send [ %s ]\n", cmd.String())
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		return err
	}
	outStr, errStr := stdout.String(), stderr.String()
	if outStr != "" || errStr != "" {
		fmt.Printf("recv [out: %s err: %s]\n", outStr, errStr)
	}
	return
}

func (a *AdbClient) getAddr() string {
	return fmt.Sprintf("%s:%d", a.IP, a.Port)
}

//clickEventCode key code event
func (a *AdbClient) eventCode(code string) error {
	cmd := exec.Command(getAdbCli(), "-s", a.getAddr(), "shell", "input", "keyevent", code)
	return cmdRun(cmd)
}

//click click screen
func (a *AdbClient) click(x, y string) error {
	cmd := exec.Command(getAdbCli(), "-s", a.getAddr(), "shell", "input", "tap", x, y)
	return cmdRun(cmd)
}

//Command
func (a *AdbClient) command(args ...string) error {
	cmd := exec.Command(getAdbCli(), args...)
	return cmdRun(cmd)
}

//Swipe swipe screen
func (a *AdbClient) swipe(StartX, StartY, EndX, EndY int) error {
	cmd := exec.Command(getAdbCli(), "-s", a.getAddr(), "shell", "input", "swipe",
		fmt.Sprintf("%d", StartX),
		fmt.Sprintf("%d", StartY),
		fmt.Sprintf("%d", EndX),
		fmt.Sprintf("%d", EndY),
	)
	return cmdRun(cmd)
}

//connect connect device
func (a *AdbClient) connect() error {
	cmd := exec.Command(getAdbCli(), "connect", a.getAddr())
	cmd.Stdout = os.Stdout
	rtn, _ := ioutil.ReadAll(os.Stdout)
	if strings.Contains(string(rtn), "connected") {
		return fmt.Errorf("conn fail:%s", rtn)
	}
	return cmdRun(cmd)
}

//disconnect disconnect device
func (a *AdbClient) disconnect() error {
	cmd := exec.Command(getAdbCli(), "disconnect", a.getAddr())
	return cmdRun(cmd)
}

//allPackage get package info
func (a *AdbClient) allPackage() error {
	cmd := exec.Command(getAdbCli(), "-s", a.getAddr(), "shell", "pm", "list", "package", "-s")
	return cmdRun(cmd)
}

//thirdPackage get third party package info
func (a *AdbClient) thirdPackage() error {
	cmd := exec.Command(getAdbCli(), "-s", a.getAddr(), "shell", "pm", "list", "package", "-3")
	return cmdRun(cmd)
}

//thirdPackage get third party package info
func (a *AdbClient) containPackage(name string) error {
	cmd := exec.Command(getAdbCli(), "-s", a.getAddr(), "shell", "pm", "list", "package", name)
	return cmdRun(cmd)
}

func (a *AdbClient) clickDialPage() error {
	cmd := exec.Command(getAdbCli(), "-s", a.getAddr(), "shell", "am", "start", "-a", "android.intent.action.DIAL")
	return cmdRun(cmd)
}

func (a *AdbClient) clickDialPhone(phone string) error {
	cmd := exec.Command(getAdbCli(), "-s", a.getAddr(), "shell", "am", "start", "-a", "android.intent.action.DIAL", "-d", fmt.Sprintf("tel:%s", phone))
	return cmdRun(cmd)
}

func (a *AdbClient) input(content string) error {
	cmd := exec.Command(getAdbCli(), "-s", a.getAddr(), "shell", "input", "text", content)
	return cmdRun(cmd)
}

//runApp
func (a *AdbClient) runApp(appPath string) error {
	cmd := exec.Command(getAdbCli(), "-s", a.getAddr(), "shell", "am", "start", "-n", appPath)
	return cmdRun(cmd)
}

//closeApp
func (a *AdbClient) closeApp(packageName string) error {
	cmd := exec.Command(getAdbCli(), "shell", "am", "force-stop", packageName)
	return cmdRun(cmd)
}
