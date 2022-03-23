package adb

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

const AdbServerPort = 5037
const AdbDaemonPort = 5555

var adbExePath string = "adb/win/adb.exe"
var adbWinApi string = "adb/win/AdbWinApi.dll"
var adbWinUsbApi string = "adb/win/AdbWinUsbApi.dll"

//getAdbCli return adb cli
func getAdbCli() string {
	if runtime.GOOS == "windows" {
		_, err := os.Stat(adbExePath)
		if err == nil {
			return adbExePath
		}
		if os.IsNotExist(err) {
			if err := downFile("https://github.com/Eminlin/adbCtrlGo/raw/main/adb/win/adb.exe", adbExePath); err != nil {
				fmt.Println(err)
				return ""
			}
			if err := downFile("https://github.com/Eminlin/adbCtrlGo/raw/main/adb/win/AdbWinApi.dll", adbWinApi); err != nil {
				fmt.Println(err)
				return ""
			}
			if err := downFile("https://github.com/Eminlin/adbCtrlGo/raw/main/adb/win/AdbWinUsbApi.dll", adbWinUsbApi); err != nil {
				fmt.Println(err)
				return ""
			}
			return adbExePath
		}
		return adbExePath
	}
	return "adb"
}

func downFile(url, path string) error {
	fmt.Printf("download file:%s \n", url)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
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
		return
	}
	return
}

//connect connect device
func (a *AdbClient) connect() error {
	cmd := exec.Command(getAdbCli(), "connect", a.getAddr())
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return err
	}
	outStr := stdout.String()
	if !strings.Contains(string(outStr), "connected") {
		a.tcpip(a.Port)
		fmt.Println("try again")
		return fmt.Errorf("conn fail:%s", outStr)
	}
	return cmdRun(cmd)
}

//disconnect disconnect device
func (a *AdbClient) disconnect() error {
	cmd := exec.Command(getAdbCli(), "disconnect", a.getAddr())
	return cmdRun(cmd)
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

func (a *AdbClient) tcpip(port int) error {
	cmd := exec.Command(getAdbCli(), "tpcip", fmt.Sprintf("%d", port))
	return cmdRun(cmd)
}

func (a *AdbClient) getAppPathByPack(packname string) (string, error) {
	cmd := exec.Command(getAdbCli(), "-s", a.getAddr(), "shell", "pm", "path", packname)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	outStr := stdout.String()
	return outStr, nil
}

//getElement
func (a *AdbClient) getElement(filename string) error {
	cmd := exec.Command(getAdbCli(), "-s", a.getAddr(), "shell", "uiautomator", "dump", "/sdcard/"+filename+".xml")
	return cmdRun(cmd)
}

func (a *AdbClient) downFile(filename, tempPath string) error {
	cmd := exec.Command(getAdbCli(), "-s", a.getAddr(), "pull", "/sdcard/"+filename+".xml", tempPath)
	return cmdRun(cmd)
}

func (a *AdbClient) getPackInfo(name string) (string, error) {
	cmd := exec.Command(getAdbCli(), "-s", a.getAddr(), "shell", "dumpsys", "package", name)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	outStr := stdout.String()
	return outStr, nil
}
