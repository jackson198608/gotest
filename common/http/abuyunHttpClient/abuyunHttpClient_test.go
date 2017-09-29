package abuyunHttpClient

import (
	"fmt"
	"net/http"
	"testing"
)

const proxyServer = ""
const proxyUser = ""
const proxyPasswd = ""

func TestAbuyun(t *testing.T) {
	var abuyun *AbuyunProxy = NewAbuyunProxy(proxyServer,
		proxyUser,
		proxyPasswd)

	t.Log("begin the test")

	if abuyun == nil {
		t.Error("create abuyun error")
	}

	targetUrl := "http://m.goumin.com/"

	var h http.Header = make(http.Header)
	h.Set("a", "1")
	statusCode, responseHeader, _, err := abuyun.SendRequest(targetUrl, h, true)
	fmt.Println(statusCode)
	fmt.Println(responseHeader)
	//fmt.Println(body)
	fmt.Println(err)

}