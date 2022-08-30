package service

import (
	"fmt"
	"score-calculate/logger"
	"strconv"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

const (
	chromeDriver = "./driver/chromedriver"
	port         = 9515
	loginPage    = "https://zjuam.zju.edu.cn/cas/login?service=http://jwbinfosys.zju.edu.cn/default2.aspx"
	searchPage   = "http://jwbinfosys.zju.edu.cn/xscj.aspx?xh="
	year         = 2022
)

var (
	service *selenium.Service
	w_b1    selenium.WebDriver
	// err     error
)

func dealerr(err error) {
	if err != nil {
		panic(err)
	}
}

func startService() *selenium.Service {
	defer func() {
		err := recover()
		if err != nil {
			logger.Logger().WithField(
				"Place", "service",
			).Error("启动服务出错")
			logger.Logger().WithField(
				"Place", "service",
			).Error(err)
		} else {
			logger.Logger().WithField(
				"Place", "service",
			).Info("启动服务")
		}
	}()
	opts := []selenium.ServiceOption{}
	//opts := []selenium.ServiceOption{
	//    selenium.StartFrameBuffer(),           // Start an X frame buffer for the browser to run in.
	//    selenium.GeckoDriver(geckoDriverPath), // Specify the path to GeckoDriver in order to use Firefox.
	//}

	// selenium.SetDebug(true)
	service, err := selenium.NewChromeDriverService(chromeDriver, port, opts...)
	dealerr(err)
	return service
	//注意这里，server关闭之后，chrome窗口也会关闭
}

func startRemote() selenium.WebDriver {
	defer func() {
		err := recover()
		if err != nil {
			logger.Logger().WithField(
				"Place", "remote",
			).Error("打开窗口出错")
			logger.Logger().WithField(
				"Place", "remote",
			).Error(err)
		} else {
			logger.Logger().WithField(
				"Place", "remote",
			).Info("打开窗口")
		}
	}()
	//链接本地的浏览器 chrome
	caps := selenium.Capabilities{
		"browserName": "chrome",
	}

	//禁止图片加载，加快渲染速度
	imagCaps := map[string]interface{}{
		"profile.managed_default_content_settings.images": 2,
	}
	chromeCaps := chrome.Capabilities{
		Prefs: imagCaps,
		Path:  "",
		Args: []string{
			"--headless", // 设置Chrome无头模式，在linux下运行，需要设置这个参数，否则会报错
			"--no-sandbox",
			"--user-agent=Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36", // 模拟user-agent，防反爬
		},
	}
	//以上是设置浏览器参数
	caps.AddChrome(chromeCaps)

	// 调起chrome浏览器
	w_b1, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	dealerr(err)
	return w_b1
}

func login(id, passwd string) {
	defer func() {
		err := recover()
		if err != nil {
			logger.Logger().WithField(
				"Place", "login",
			).Error("登录出错")
			logger.Logger().WithField(
				"Place", "login",
			).Error(err)
		} else {
			logger.Logger().WithField(
				"Place", "login",
			).Info("登录")
		}
	}()
	if err := w_b1.Get(loginPage); err != nil {
		panic(err)
	}

	// 填写信息
	elem, err := w_b1.FindElement(selenium.ByCSSSelector, "#username")
	dealerr(err)
	elem.SendKeys(id)

	elem, err = w_b1.FindElement(selenium.ByCSSSelector, "#password")
	dealerr(err)
	elem.SendKeys(passwd)

	elem, err = w_b1.FindElement(selenium.ByCSSSelector, "#dl")
	dealerr(err)
	elem.Click()

}

func searchSem(myPage string, year int) ([]string, []int, []float32, []float32) {
	defer func() {
		err := recover()
		if err != nil {
			logger.Logger().WithField(
				"Place", "search",
			).Error("搜索信息出错")
			logger.Logger().WithField(
				"Place", "search",
			).Error(err)
		} else {
			logger.Logger().WithField(
				"Place", "search",
			).Info("搜索窗口")
		}
	}()
	err := w_b1.Get(myPage)
	dealerr(err)

	choiceStr := fmt.Sprintf("//select[@id='ddlXN']/option[@value='%d-%d']", year-1, year)
	elem, err := w_b1.FindElement(selenium.ByXPATH, choiceStr)
	dealerr(err)
	elem.Click()

	// elem, err = w_b1.FindElement(selenium.ByCSSSelector, "#Button5")
	// dealerr(err)
	w_b1.ExecuteScript("el = document.getElementById('Button5'); el.click()", nil)

	var (
		classNames []string
		scores     []int
		credits    []float32
		gpas       []float32
	)

	elems, err := w_b1.FindElements(selenium.ByCSSSelector, "#DataGrid1 tr")
	dealerr(err)
	for item, elem := range elems {
		if item == 0 {
			continue
		}
		parts, _ := elem.FindElements(selenium.ByCSSSelector, "td")
		className, _ := parts[1].Text()
		text, _ := parts[2].Text()
		score, _ := strconv.ParseInt(text, 10, 8)
		text, _ = parts[3].Text()
		credit, _ := strconv.ParseFloat(text, 32)
		text, _ = parts[4].Text()
		gpa, _ := strconv.ParseFloat(text, 32)

		classNames = append(classNames, className)
		scores = append(scores, int(score))
		credits = append(credits, float32(credit))
		gpas = append(gpas, float32(gpa))
	}

	return classNames, scores, credits, gpas
}

func GetAllMess(id, passwd string) map[string]SemResult {
	result := make(map[string]SemResult)

	myPage := searchPage + id

	w_b1 = startRemote()
	defer w_b1.Quit()

	login(id, passwd)

	for yi := 0; yi < 6; yi++ {
		a, b, c, d := searchSem(myPage, year-yi)
		yearDur := fmt.Sprintf("%d-%d", year-yi-1, year-yi)
		if a == nil {
			break
		}
		result[yearDur] = SemResult{
			ClassNames: a,
			Scores:     b,
			Credits:    c,
			Gpas:       d,
		}
	}
	return result
}
