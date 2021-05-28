package youdao

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	signType = "v3"
	ApiUrl   = "https://openapi.youdao.com/api"
)

//TransField 翻译结果结构体
type TransField struct {
	UsPhonetic string   `json:"us-phonetic"`
	Phonetic   string   `json:"phonetic"`
	UkPhonetic string   `json:"uk-phonetic"`
	UkSpeech   string   `json:"uk-speech"`
	UsSpeech   string   `json:"us-speech"`
	Explains   []string `json:"explains"`
}

//TransWeb 延伸释义结构体
type TransWeb struct {
	Key   string   `json:"key"`
	Value []string `json:"value"`
}

//Config 配置文件
type Config struct {
	AppKey    string `json:"AppKey"`
	AppSecret string `json:"AppSecret"`
}

type Context struct {
	fromLang string
	toLang   string
	q        string
}
type Client struct {
	config *Config
}
type TextTranslationReq struct {
	fromLang                 string
	toLang                   string
	q                        string
	signToStr, salt, curTime string
}
type TextTranslationResp struct {
	ErrorCode    string                 `json:"errorCode"`
	Query        string                 `json:"query"`
	Translation  []string               `json:"translation"`
	Basic        TransField             `json:"basic"`
	Web          []TransWeb             `json:"web,omitempty"`
	Language     string                 `json:"l"`
	Dict         map[string]interface{} `json:"dict,omitempty"`
	WebDict      map[string]interface{} `json:"webdict,omitempty"`
	TSpeakUrl    string                 `json:"tSpeakUrl,omitempty"`
	SpeakUrl     string                 `json:"speakUrl,omitempty"`
	ReturnPhrase []string               `json:"returnPhrase,omitempty"`
}

//DictionaryRespJson 查询回应报文结构体
type DictionaryRespJson struct {
	ErrorCode    string                 `json:"errorCode"`
	Query        string                 `json:"query"`
	Translation  []string               `json:"translation"`
	Basic        TransField             `json:"basic"`
	Web          []TransWeb             `json:"web,omitempty"`
	Language     string                 `json:"l"`
	Dict         map[string]interface{} `json:"dict,omitempty"`
	WebDict      map[string]interface{} `json:"webdict,omitempty"`
	TSpeakUrl    string                 `json:"tSpeakUrl,omitempty"`
	SpeakUrl     string                 `json:"speakUrl,omitempty"`
	ReturnPhrase []string               `json:"returnPhrase,omitempty"`
}

//TextTranslation 调用API，查询翻译结果
func (client *Client) TextTranslation(ctx context.Context, req TextTranslationReq) (resp TextTranslationResp, err error) {

	data := make(url.Values, 0)
	data["q"] = []string{req.q}
	data["from"] = []string{req.fromLang}
	data["to"] = []string{req.toLang}
	data["appKey"] = []string{client.config.AppKey}
	data["salt"] = []string{req.salt}
	data["sign"] = []string{req.signToStr}
	data["signType"] = []string{signType}
	data["curtime"] = []string{req.curTime}

	respon, err := http.PostForm(ApiUrl, data)
	if err != nil {
		return TextTranslationResp{}, fmt.Errorf("%w", err)
	}
	defer respon.Body.Close()
	body, err := ioutil.ReadAll(respon.Body)
	if err != nil {
		return TextTranslationResp{}, fmt.Errorf("%w", err)
	}

	err = json.Unmarshal(body, &resp)
	if err != nil {
		return TextTranslationResp{}, fmt.Errorf("%w", err)
	}
	return resp, nil
}

//newClient 初始化
func NewClient(config *Config) Client {
	var client Client
	client.config = config
	return client
}
func NewTextTranslateReq(fromLang, toLang, q, signToStr, salt, curTime string) (ttr TextTranslationReq) {
	ttr.toLang = fromLang
	ttr.fromLang = fromLang
	ttr.q = q
	ttr.salt = salt
	ttr.signToStr = signToStr
	ttr.curTime = curTime
	return
}
