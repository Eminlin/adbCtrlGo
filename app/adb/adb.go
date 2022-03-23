package adb

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
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

var dumpFilePath string = "/sdcard/dump.xml"

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
func (a *AdbClient) cmdRun(cmd *exec.Cmd) (err error) {
	if a.Debug {
		fmt.Printf("send [ %s ]\n", cmd.String())
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		return err
	}
	outStr := stdout.String()
	errStr := stderr.String()
	if a.Debug {
		fmt.Printf("recv [out:%s]\n", outStr)
	}
	if errStr != "" {
		return errors.New(errStr)
	}
	return
}

func (a *AdbClient) cmdRunRtn(cmd *exec.Cmd) (string, error) {
	if a.Debug {
		fmt.Printf("send [ %s ]\n", cmd.String())
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	outStr := stdout.String()
	errStr := stderr.String()
	if errStr != "" {
		return outStr, errors.New(errStr)
	}
	if a.Debug {
		fmt.Printf("recv [out: %s]\n", outStr)
	}
	return outStr, nil
}

type realtime struct {
	out string
}

func (a *AdbClient) cmdRunRealTime(cmd *exec.Cmd, out *realtime) error {
	if a.Debug {
		fmt.Printf("send [ %s ]\n", cmd.String())
	}
	stderr, _ := cmd.StderrPipe()
	stdout, _ := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil {
		return err
	}
	logScan := bufio.NewScanner(stdout)
	for logScan.Scan() {
		if a.Debug {
			fmt.Printf("recv [ %s ]", logScan.Text())
		}
		out.out += logScan.Text() + "\n"
	}

	errBuf := bytes.NewBufferString("")
	scan := bufio.NewScanner(stderr)
	for scan.Scan() {
		s := scan.Text()
		log.Println("build error: ", s)
		errBuf.WriteString(s)
		errBuf.WriteString("\n")
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	if !cmd.ProcessState.Success() {
		return errors.New(errBuf.String())
	}
	return nil
}

//connect connect device
func (a *AdbClient) connect() error {
	for i := a.ConnectMaxTryTimes; i == 0; i-- {
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
			if err := a.tcpip(a.Port); err != nil {
				return err
			}
			fmt.Printf("conn fail:%s \n", outStr)
			fmt.Printf("try again(%d)...", i)
			continue
		}
		if err := a.cmdRun(cmd); err != nil {
			return err
		}
	}
	return nil
}

//disconnect disconnect device
func (a *AdbClient) disconnect() error {
	cmd := exec.Command(getAdbCli(), "disconnect", a.getAddr())
	return a.cmdRun(cmd)
}

//disconnect disconnect device
func (a *AdbClient) getAdbVersion() string {
	cmd := exec.Command(getAdbCli(), "version")
	out, err := a.cmdRunRtn(cmd)
	if err != nil {
		return "unknown"
	}
	return out
}

func (a *AdbClient) getAddr() string {
	return fmt.Sprintf("%s:%d", a.IP, a.Port)
}

//clickEventCode key code event
func (a *AdbClient) eventCode(code string) error {
	cmd := exec.Command(getAdbCli(), "-s", a.getAddr(), "shell", "input", "keyevent", code)
	return a.cmdRun(cmd)
}

//click click screen
func (a *AdbClient) tap(x, y string) error {
	cmd := exec.Command(getAdbCli(), "-s", a.getAddr(), "shell", "input", "tap", x, y)
	return a.cmdRun(cmd)
}

//Command
func (a *AdbClient) command(args ...string) error {
	cmd := exec.Command(getAdbCli(), args...)
	return a.cmdRun(cmd)
}

//Swipe swipe screen
func (a *AdbClient) swipe(StartX, StartY, EndX, EndY int) error {
	cmd := exec.Command(getAdbCli(), "-s", a.getAddr(), "shell", "input", "swipe",
		fmt.Sprintf("%d", StartX),
		fmt.Sprintf("%d", StartY),
		fmt.Sprintf("%d", EndX),
		fmt.Sprintf("%d", EndY),
	)
	return a.cmdRun(cmd)
}

//allPackage get package info
func (a *AdbClient) allPackage() error {
	cmd := exec.Command(getAdbCli(), "-s", a.getAddr(), "shell", "pm", "list", "package", "-s")
	return a.cmdRun(cmd)
}

//thirdPackage get third party package info
func (a *AdbClient) thirdPackage() error {
	cmd := exec.Command(getAdbCli(), "-s", a.getAddr(), "shell", "pm", "list", "package", "-3")
	return a.cmdRun(cmd)
}

//thirdPackage get third party package info
func (a *AdbClient) containPackage(name string) error {
	cmd := exec.Command(getAdbCli(), "-s", a.getAddr(), "shell", "pm", "list", "package", name)
	return a.cmdRun(cmd)
}

func (a *AdbClient) tapDialPage() error {
	cmd := exec.Command(getAdbCli(), "-s", a.getAddr(), "shell", "am", "start", "-a", "android.intent.action.DIAL")
	return a.cmdRun(cmd)
}

func (a *AdbClient) tapDialPhone(phone string) error {
	cmd := exec.Command(getAdbCli(), "-s", a.getAddr(), "shell", "am", "start", "-a", "android.intent.action.DIAL", "-d", fmt.Sprintf("tel:%s", phone))
	return a.cmdRun(cmd)
}

func (a *AdbClient) input(content string) error {
	cmd := exec.Command(getAdbCli(), "-s", a.getAddr(), "shell", "input", "text", content)
	return a.cmdRun(cmd)
}

//runApp
func (a *AdbClient) runApp(appPath string) error {
	cmd := exec.Command(getAdbCli(), "-s", a.getAddr(), "shell", "am", "start", "-n", appPath)
	return a.cmdRun(cmd)
}

//closeApp
func (a *AdbClient) forceStopApp(packageName string) error {
	cmd := exec.Command(getAdbCli(), "shell", "am", "force-stop", packageName)
	return a.cmdRun(cmd)
}

func (a *AdbClient) tcpip(port int) error {
	cmd := exec.Command(getAdbCli(), "tpcip", fmt.Sprintf("%d", port))
	return a.cmdRun(cmd)
}

func (a *AdbClient) getAppPathByPack(packname string) (string, error) {
	cmd := exec.Command(getAdbCli(), "-s", a.getAddr(), "shell", "pm", "path", packname)
	return a.cmdRunRtn(cmd)
}

//getElement
func (a *AdbClient) getUiautomatorElement(filename string) error {
	cmd := exec.Command(getAdbCli(), "-s", a.getAddr(), "shell", "uiautomator", "dump", dumpFilePath)
	return a.cmdRun(cmd)
}

func (a *AdbClient) downFile(filePath, tempPath string) error {
	cmd := exec.Command(getAdbCli(), "-s", a.getAddr(), "pull", filePath, tempPath)
	return a.cmdRun(cmd)
}

func (a *AdbClient) push(localPath, remotePath string) error {
	cmd := exec.Command(getAdbCli(), "-s", a.getAddr(), "pull", localPath, remotePath)
	return a.cmdRun(cmd)
}

func (a *AdbClient) getPackInfo(name string) (string, error) {
	cmd := exec.Command(getAdbCli(), "-s", a.getAddr(), "shell", "dumpsys", "package", name)
	return a.cmdRunRtn(cmd)
}

func (a *AdbClient) killserver() error {
	cmd := exec.Command(getAdbCli(), "-s", a.getAddr(), "kill-server")
	return a.cmdRun(cmd)
}

func (a *AdbClient) install(apkPath string) error {
	var apkName string
	if strings.Contains(apkPath, "/") {
		s := strings.Split(apkPath, "/")
		apkName = s[len(s)-1]
	} else {
		apkName = apkPath
	}
	if err := a.push(apkPath, "/data/local/tmp"); err != nil {
		return err
	}
	cmd := exec.Command(getAdbCli(), "-s", a.getAddr(), "install", "/data/local/tmp/"+apkName)
	return a.cmdRun(cmd)
}

func (a *AdbClient) ping(url string) (string, error) {
	cmd := exec.Command(getAdbCli(), "-s", a.getAddr(), "shell", "ping", "-c", "4", url)
	out := realtime{}
	err := a.cmdRunRealTime(cmd, &out)
	return out.out, err
}

func (a *AdbClient) getPhoneModel() (string, error) {
	cmd := exec.Command(getAdbCli(), "-s", a.getAddr(), "shell", "getprop", "ro.product.model")
	return a.cmdRunRtn(cmd)
}

func (a *AdbClient) getBatterryState() (string, error) {
	cmd := exec.Command(getAdbCli(), "-s", a.getAddr(), "shell", "dumpsys", "battery")
	return a.cmdRunRtn(cmd)
}

func (a *AdbClient) getScreenSize() (string, error) {
	cmd := exec.Command(getAdbCli(), "-s", a.getAddr(), "shell", "vm", "size")
	return a.cmdRunRtn(cmd)
}

func (a *AdbClient) getScreenDensity() (string, error) {
	cmd := exec.Command(getAdbCli(), "-s", a.getAddr(), "shell", "vm", "density")
	return a.cmdRunRtn(cmd)
}

func (a *AdbClient) getAndroidID() (string, error) {
	cmd := exec.Command(getAdbCli(), "-s", a.getAddr(), "shell", "settings", "get", "secure", "android_id")
	return a.cmdRunRtn(cmd)
}

func (a *AdbClient) getAndroidVersion() (string, error) {
	cmd := exec.Command(getAdbCli(), "-s", a.getAddr(), "shell", "getprop", "ro.build.version.release")
	return a.cmdRunRtn(cmd)
}

func (a *AdbClient) getCPUInfo() (string, error) {
	cmd := exec.Command(getAdbCli(), "-s", a.getAddr(), "shell", "cat", "/proc/cpuinfo")
	return a.cmdRunRtn(cmd)
}

func (a *AdbClient) getMemoryInfo() (string, error) {
	cmd := exec.Command(getAdbCli(), "-s", a.getAddr(), "shell", "cat", "/proc/meminfo")
	return a.cmdRunRtn(cmd)
}

func (a *AdbClient) getPhoneBrand() (string, error) {
	cmd := exec.Command(getAdbCli(), "-s", a.getAddr(), "shell", "getprop", "ro.product.brand")
	return a.cmdRunRtn(cmd)
}

func (a *AdbClient) reboot() (string, error) {
	cmd := exec.Command(getAdbCli(), "-s", a.getAddr(), "reboot")
	return a.cmdRunRtn(cmd)
}

func (a *AdbClient) screenshotPng(outPath string) error {
	cmd := exec.Command(getAdbCli(), "-s", a.getAddr(), "exec-out", "screencap", "/sdcard/screenshot.png")
	if err := a.downFile("/sdcard/screenshot.png", outPath); err != nil {
		return err
	}
	return a.cmdRun(cmd)
}
