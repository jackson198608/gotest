package SaveHtml

import (
	"errors"
	"fmt"
	"github.com/donnie4w/go-logger/logger"
	"github.com/jackson198608/goProject/common/http/abuyunHttpClient"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
)

type HtmlInfo struct {
	id        int
	url       string
	queueName string
	saveDir   string
	abuyun    *abuyunHttpClient.AbuyunProxy
}

func NewHtml(logLevel int, queueName string, id int, url string, taskNewArgs []string, abuyun *abuyunHttpClient.AbuyunProxy) *HtmlInfo {
	logger.SetLevel(logger.LEVEL(logLevel))
	e := new(HtmlInfo)
	e.id = id
	e.url = url
	e.queueName = queueName
	e.saveDir = taskNewArgs[0] // 0:saveDir
	e.abuyun = abuyun
	return e
}

//get content by url
func (e *HtmlInfo) CreateHtmlByUrl() error {
	statusCode, _, body, err := e.changeIpByAbuyun()
	if err != nil {
		logger.Error("change ip abuyun error", err)
		return errors.New("change ip abuyun error")
	}
	fmt.Println(statusCode, e.id)
	if statusCode == 200 {
		urlname := e.saveFileName()
		status := e.saveContentToHtml(urlname, body)
		if status == true {
			logger.Info("save content to html: ", urlname)
			return nil
		}
		return errors.New("save content html error")
	} else {
		return errors.New("get html error")
	}
}

// change ip by abuyun
func (e *HtmlInfo) changeIpByAbuyun() (int, *http.Header, string, error) {
	logger.Info("begin the test", e.id)

	if e.abuyun == nil {
		logger.Error("create abuyun error")
	}
	var h http.Header = make(http.Header)
	statusCode, responseHeader, body, err := e.abuyun.SendRequest(e.url, h, true)
	return statusCode, responseHeader, body, err
}

// create filename
func (e *HtmlInfo) saveFileName() string {
	filename := ""
	dir := ""
	if e.id < 1000 {
		dir = ""
	} else {
		n4 := e.id % 10               //个位数
		n3 := (e.id - n4) % 100       //十位数
		n2 := (e.id - n4 - n3) % 1000 //百位数
		dir = strconv.Itoa(n2/100) + "/" + strconv.Itoa(n3/10) + "/" + strconv.Itoa(n4) + "/"
	}
	urlstr := strings.Split(e.url, "/")
	strlen := len(urlstr)
	if strlen >= 1 {
		urlstrReal := strings.Split(urlstr[strlen-1], "?")
		filename = dir + urlstrReal[0]
	}
	return filename
}

// save content to html file
func (e *HtmlInfo) saveContentToHtml(urlname string, content string) bool {
	var filename = e.saveDir + urlname
	var f *os.File
	var err1 error
	if checkFileIsExist(filename) { //如果文件存在
		f, err1 = os.OpenFile(filename, os.O_APPEND, 0666) //打开文件
		fmt.Println("文件存在")
	} else {
		f, err1 = os.Create(filename) //创建文件
		fmt.Println("文件不存在")
	}
	fullDirPath := path.Dir(filename)
	err1 = os.MkdirAll(fullDirPath, 0777)
	check(err1)
	f, err1 = os.Create(filename) //创建文件
	check(err1)
	n, err1 := f.WriteString(content) //写入文件(字符串)
	check(err1)
	fmt.Printf("写入 %d 个字节n", n)
	return true
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

/**
 * 判断文件是否存在  存在返回 true 不存在返回false
 */
func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}
