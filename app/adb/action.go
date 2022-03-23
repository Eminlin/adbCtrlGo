package adb

type AdbClient struct {
	IP   string
	Port int
}

func NewAdbClient(ip string, port int) (*AdbClient, error) {
	client := &AdbClient{IP: ip, Port: port}
	if err := client.connect(); err != nil {
		return client, err
	}
	return client, nil
}

//ClickHome click home menu
func (a *AdbClient) ClickHome() error {
	return a.eventCode("3")
}

//ClickBack click back menu
func (a *AdbClient) ClickBack() error {
	return a.eventCode("4")
}

//Click click screen
func (a *AdbClient) Click(x, y string) error {
	return a.click(x, y)
}

//Command exec some command
func (a *AdbClient) Command(args ...string) error {
	return a.command(args...)
}

func (a *AdbClient) EventCode(code string) error {
	return a.eventCode(code)
}

func (a *AdbClient) Swipe(startX, startY, endX, endY int) error {
	return a.swipe(startX, startY, endX, endY)
}

func (a *AdbClient) Power() error {
	return a.eventCode("26")
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

func (a *AdbClient) ClickDialPage() error {
	return a.clickDialPage()
}

func (a *AdbClient) ClickDialPhone(phone string) error {
	return a.clickDialPhone(phone)
}

func (a *AdbClient) Input(content string) error {
	return a.input(content)
}

func (a *AdbClient) RunApp(appPath string) error {
	return a.runApp(appPath)
}

func (a *AdbClient) CloseApp(appPath string) error {
	return a.closeApp(appPath)
}

func (a *AdbClient) GetAppPath(packname string) (string, error) {
	return a.getAppPathByPack(packname)
}

func (a *AdbClient) GetElement(filename string) error {
	return a.getElement(filename)
}

func (a *AdbClient) Down(file string) error {
	return a.downFile(file, "temp")
}

func (a *AdbClient) GetAppInfo(packname string) (string, error) {
	return a.getPackInfo(packname)
}

func (a *AdbClient) LowLight() error {
	return a.eventCode("220")
}
