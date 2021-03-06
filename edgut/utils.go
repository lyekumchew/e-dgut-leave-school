package edgut

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/lyekumchew/e-dgut-leave-school/common"
	"github.com/lyekumchew/e-dgut-leave-school/config"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
	"time"
)

const (
	homeUrl                  = "http://e.dgut.edu.cn"
	ibpstestUrl              = "https://cas.dgut.edu.cn/home/Oauth/getToken/appid/ibpstest/state/home"
	studentLeaveOnLoadDaoUrl = "http://219.222.186.78:17750/api/studentLeaveOnLoadDao"
	getBoDataUrl             = "http://e.dgut.edu.cn/api/adminhome/getBoData"
	getFormDataUrl           = "http://e.dgut.edu.cn/ibps/business/v3/bpm/instance/getFormData"
	getUserInfoUrl           = "http://e.dgut.edu.cn/api/cas/getUserInfo"
	applyUrl                 = "http://e.dgut.edu.cn/ibps/business/v3/bpm/instance/start"
	defId                    = "758369743466921984"
	scUrl                    = "https://sc.ftqq.com/"
)

type EDGUTClient struct {
	Config config.Config
	token  string
	Data   Data
}

type Header struct {
	Key   string
	Value string
}

type Data struct {
	Parameters [3]struct {
		Key   interface{} `json:"key"`
		Value interface{} `json:"value"`
	} `json:"parameters"`
}

type Value struct {
	XueHao              string `json:"xuehao"`
	ShenPiRen           string `json:"shenPiRen"`
	BaiMingDanQuanXian  string `json:"baiMingDanQuanXian"`
	FanXiaoLuXian       string `json:"fanXiaoLuXian"`
	FanXiaoChengZuoJTGJ string `json:"fanXiaoChengZuoJTGJ"`
	LiXiaoLuXian        string `json:"liXiaoLuXian"`
	LiXiaoChengZuoJTGJ  string `json:"liXiaoChengZuoJTGJ"`
	JiaTingZhuZhi       string `json:"jiaTingZhuZhi"`
	JiaChangDianHua     string `json:"jiaChangDianHua"`
	QingJiaYuanYin      string `json:"qingJiaYuanYin"`
	LiXiaoMuDiDi        string `json:"liXiaoMuDiDi"`
	QingJiaLeiXing      string `json:"qingJiaLeiXing"`
	QingJiaTianShu      int    `json:"qingJiaTianShu`
	FanXiaoShiJian      string `json:"fanXiaoShiJian"`
	LiXiaoShiJian       string `json:"liXiaoShiJian"`
	LianXiDianHua       string `json:"lianXiDianHua"`
	BanJi               string `json:"banJi"`
	ZhuanYe             string `json:"zhuanYe"`
	Id                  string `json:"id"`
}

var client http.Client

func init() {
	// http.client init
	jar, _ := cookiejar.New(nil)
	client = http.Client{Jar: jar}
}

func (e *EDGUTClient) Login() (err error) {
	// fetch the xss token
	resp, err := client.Get(ibpstestUrl)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	contents, _ := ioutil.ReadAll(resp.Body)
	re := regexp.MustCompile(`var token = "(.*?)";`)
	res := re.FindAllStringSubmatch(string(contents), -1)
	xssToken := res[0][1]
	if xssToken == "" {
		return errors.New("cant not fetch the xss token")
	}

	// login params
	params := url.Values{}
	params.Set("username", e.Config.Username)
	params.Set("password", e.Config.Password)
	params.Set("__token__", xssToken)

	// post -> homeUrl
	req, _ := http.NewRequest("POST", ibpstestUrl, strings.NewReader(params.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp2, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp2.Body.Close()
	contents2, _ := ioutil.ReadAll(resp2.Body)

	if strings.Contains(string(contents2), "通过") {
	} else {
		return errors.New("cas login error, msg: " + string(contents2))
	}

	// access token
	re = regexp.MustCompile(`"info":"(.*?)"}`)
	res = re.FindAllStringSubmatch(string(contents2), -1)
	redirectURl := res[0][1]
	redirectURl = strings.ReplaceAll(redirectURl, "\\", "")
	resp3, err := client.Get(redirectURl)
	defer resp3.Body.Close()
	if err != nil {
		return err
	}
	re = regexp.MustCompile(`access_token=(.*?)$`)
	res = re.FindAllStringSubmatch(resp3.Request.URL.String(), -1)
	e.token = strings.Split(res[0][1], "&")[0]
	if e.token == "" {
		return errors.New("fetch the access token error")
	}

	return nil
}

func fetch(method, _url, token string, headers ...Header) (s string, err error) {
	req, _ := http.NewRequest(method, _url, strings.NewReader(url.Values{}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("x-authorization-access_token", token)
	if len(headers) > 0 {
		for _, v := range headers {
			req.Header.Set(v.Key, v.Value)
		}
	}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}
	contents, _ := ioutil.ReadAll(resp.Body)
	return string(contents), nil
}

func (e *EDGUTClient) Do() error {
	contents, _ := fetch("GET", getUserInfoUrl, e.token)
	orgId := gjson.Get(contents, "info.orgs.id").String()

	contents, _ = fetch("GET", getFormDataUrl+"?defId="+defId, e.token)
	re := regexp.MustCompile(`code=(.*?)&`)
	code := re.FindAllStringSubmatch(contents, -1)[0][1]

	contents, _ = fetch("GET", getBoDataUrl+"?code="+code+"&field=xue_yuan_&value="+orgId, e.token)
	approvers := gjson.Get(contents, "info.0.shen_pi_ren_").String()

	contents, _ = fetch("GET", studentLeaveOnLoadDaoUrl, e.token, Header{Key: "Origin", Value: homeUrl})
	major := gjson.Get(contents, "data.dataResult.major").String()
	class := gjson.Get(contents, "data.dataResult.classes").String()

	today := time.Now().Format("2006-1-2")

	// data
	e.Data.Parameters[0].Key = "defId"
	e.Data.Parameters[0].Value = defId
	e.Data.Parameters[1].Key = "version"
	e.Data.Parameters[1].Value = "0"
	e.Data.Parameters[2].Key = "data"
	value := Value{
		XueHao:              e.Config.Username,
		ShenPiRen:           approvers,
		BaiMingDanQuanXian:  "C",
		FanXiaoLuXian:       e.Config.ReturnRoute,
		FanXiaoChengZuoJTGJ: e.Config.ReturnRtransportation,
		LiXiaoLuXian:        e.Config.LeaveRoute,
		LiXiaoChengZuoJTGJ:  e.Config.LeaveTransportation,
		JiaTingZhuZhi:       "{\"street\":\"\",\"province\":\"44\",\"city\":\"4401\",\"district\":\"440111\"}",
		JiaChangDianHua:     e.Config.ParentsPhone,
		QingJiaYuanYin:      e.Config.ReasonDetails,
		LiXiaoMuDiDi:        "{\"street\":\"\",\"province\":\"44\",\"city\":\"4401\",\"district\":\"440111\"}",
		QingJiaLeiXing:      e.Config.LeaveReason,
		QingJiaTianShu:      0,
		FanXiaoShiJian:      today,
		LiXiaoShiJian:       today,
		LianXiDianHua:       e.Config.Contact,
		BanJi:               class,
		ZhuanYe:             major,
		Id:                  "",
	}
	j, _ := json.Marshal(value)
	e.Data.Parameters[2].Value = string(j)
	j, _ = json.Marshal(e.Data)

	req, _ := http.NewRequest("POST", applyUrl, bytes.NewBuffer(j))
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("x-authorization-access_token", e.token)
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	res, _ := ioutil.ReadAll(resp.Body)

	if strings.Contains(string(res), "流程启动成功") {
		common.SCMsg("流程启动成功", "", e.Config.SCKey)
		common.Logger("流程启动成功", 0)
	} else {
		common.SCMsg("流程启动失败", string(res), e.Config.SCKey)
		common.Logger("流程启动失败"+string(res), 1)
		return errors.New("流程启动失败: " + string(res))
	}

	return nil
}
