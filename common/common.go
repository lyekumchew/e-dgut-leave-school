package common

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

const scUrl = "https://sc.ftqq.com/"

func Logger(m string, level int) {
	color := []string{"[SUCCESS]", "[ERROR]", "[INFO]"}
	header := []string{"\u001B[32;1m", "\u001B[31;1m", "\u001B[36;1m"}
	loc, _ := time.LoadLocation("Asia/Shanghai")
	fmt.Println(header[level], color[level], time.Now().In(loc).Format("2006-01-02 15:04:05"), m, "\033[0m")
}

func SCMsg(text, desp, sckey string) {
	client := http.Client{}
	resp, err := client.Get(scUrl + sckey + ".send?text=" + url.QueryEscape(text) + "&desp=" + url.QueryEscape(desp))
	//resp, err := http.Get()
	defer resp.Body.Close()
	if err != nil && resp.Status != string(http.StatusOK) {
		Logger("Server 酱发送失败", 1)
		Logger(resp.Request.URL.String(), 2)
	}
	contents, _ := ioutil.ReadAll(resp.Body)
	if res, _ := regexp.MatchString("success", string(contents)); res != true {
		Logger("Server 酱发送失败", 1)
		Logger(string(contents), 2)
	}
}
