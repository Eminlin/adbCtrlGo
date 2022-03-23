package adb

type AdbClient struct {
	IP                 string
	Port               int
	Debug              bool
	ConnectMaxTryTimes int
}

func NewAdbClient(client AdbClient) (*AdbClient, error) {
	if err := client.connect(); err != nil {
		if err.Error() == "exec: already started" {
			return &client, nil
		}
		return &client, err
	}
	return &client, nil
}

//Tap Tap screen
func (a *AdbClient) TapScreen(x, y string) error {
	return a.tap(x, y)
}

//Command exec some command
func (a *AdbClient) Command(args ...string) error {
	return a.command(args...)
}

func (a *AdbClient) GetAdbVersion() string {
	return a.getAdbVersion()
}

func (a *AdbClient) EventCode(code string) error {
	return a.eventCode(code)
}

func (a *AdbClient) Swipe(startX, startY, endX, endY int) error {
	return a.swipe(startX, startY, endX, endY)
}

//TapHome tap home menu
func (a *AdbClient) TapHome() error {
	return a.eventCode("3")
}

//TapBack tap back menu
func (a *AdbClient) TapBack() error {
	return a.eventCode("4")
}

func (a *AdbClient) VolumeUp() error {
	return a.eventCode("24")
}

func (a *AdbClient) VolumeDown() error {
	return a.eventCode("25")
}

func (a *AdbClient) PressPower() error {
	return a.eventCode("26")
}

func (a *AdbClient) TapPhotoApp() error {
	return a.eventCode("27")
}

func (a *AdbClient) TapBrowser() error {
	return a.eventCode("64")
}

func (a *AdbClient) TapMenu() error {
	return a.eventCode("82")
}

func (a *AdbClient) TapPlayOrPause() error {
	return a.eventCode("86")
}

func (a *AdbClient) TapNextPlay() error {
	return a.eventCode("87")
}

func (a *AdbClient) TapPreviousPlay() error {
	return a.eventCode("88")
}

func (a *AdbClient) SilentMode() error {
	return a.eventCode("164")
}

func (a *AdbClient) LowerLight() error {
	return a.eventCode("220")
}

func (a *AdbClient) HigherLight() error {
	return a.eventCode("221")
}

func (a *AdbClient) SystemSleep() error {
	return a.eventCode("223")
}

func (a *AdbClient) WakesScreen() error {
	return a.eventCode("224")
}

func (a *AdbClient) GetAllPackage() error {
	return a.allPackage()
}

func (a *AdbClient) GetThirdPartPackage() error {
	return a.thirdPackage()
}

func (a *AdbClient) GetPackageByName(name string) error {
	return a.containPackage(name)
}

func (a *AdbClient) TapDialPage() error {
	return a.tapDialPage()
}

func (a *AdbClient) TapDialPhone(phone string) error {
	return a.tapDialPhone(phone)
}

func (a *AdbClient) Input(content string) error {
	return a.input(content)
}

func (a *AdbClient) RunApp(appPath string) error {
	return a.runApp(appPath)
}

func (a *AdbClient) ForceStopApp(appPath string) error {
	return a.forceStopApp(appPath)
}

func (a *AdbClient) GetAppPath(packname string) (string, error) {
	return a.getAppPathByPack(packname)
}

func (a *AdbClient) GetUiautomatorElement(filename string) error {
	return a.getUiautomatorElement(filename)
}

func (a *AdbClient) Down(filePath, tempPath string) error {
	return a.downFile(filePath, tempPath)
}

func (a *AdbClient) GetAppInfo(packname string) (string, error) {
	return a.getPackInfo(packname)
}

func (a *AdbClient) Disconnect() error {
	return a.disconnect()
}

func (a *AdbClient) KillServer() error {
	return a.killserver()
}

func (a *AdbClient) Install(apkPath string) error {
	return a.install(apkPath)
}

func (a *AdbClient) Ping(url string) (string, error) {
	return a.ping(url)
}

func (a *AdbClient) GetPhoneModel() (string, error) {
	return a.getPhoneModel()
}

func (a *AdbClient) GetBatterryState() (string, error) {
	return a.getBatterryState()
}

func (a *AdbClient) GetScreenSize() (string, error) {
	return a.getScreenSize()
}

func (a *AdbClient) GetScreenDensity() (string, error) {
	return a.getScreenDensity()
}

func (a *AdbClient) GetAndroidID() (string, error) {
	return a.getAndroidID()
}

func (a *AdbClient) GetAndroidVersion() (string, error) {
	return a.getAndroidVersion()
}

func (a *AdbClient) GetCPUInfo() (string, error) {
	return a.getCPUInfo()
}

func (a *AdbClient) GetMemoryInfo() (string, error) {
	return a.getMemoryInfo()
}

func (a *AdbClient) GetPhoneBrand() (string, error) {
	return a.getPhoneBrand()
}

func (a *AdbClient) Reboot() (string, error) {
	return a.reboot()
}

func (a *AdbClient) ScreenshotPNG(outPath string) error {
	return a.screenshotPng(outPath)
}
