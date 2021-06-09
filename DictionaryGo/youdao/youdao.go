package youdao

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	signType = "v3"
	apiUrl   = "https://openapi.youdao.com/api"
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
	Client    *http.Client
}
type Client struct {
	config Config
}
type TextTranslationReq struct {
	FromLang string //源语言
	ToLang   string //目标语言
	Q        string //待翻译文本
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

//TextTranslation 调用API，查询翻译结果
//文档地址 https://ai.youdao.com/DOCSIRMA/html/%E8%87%AA%E7%84%B6%E8%AF%AD%E8%A8%80%E7%BF%BB%E8%AF%91/API%E6%96%87%E6%A1%A3/%E6%96%87%E6%9C%AC%E7%BF%BB%E8%AF%91%E6%9C%8D%E5%8A%A1/%E6%96%87%E6%9C%AC%E7%BF%BB%E8%AF%91%E6%9C%8D%E5%8A%A1-API%E6%96%87%E6%A1%A3.html
func (client *Client) TextTranslation(ctx context.Context, req TextTranslationReq) (resp TextTranslationResp, err error) {
	signToStr, salt, curTime := getSign(req.Q, &client.config)
	data := make(url.Values, 0)
	data["q"] = []string{req.Q}
	data["from"] = []string{req.FromLang}
	data["to"] = []string{req.ToLang}
	data["appKey"] = []string{client.config.AppKey}
	data["salt"] = []string{salt}
	data["sign"] = []string{signToStr}
	data["signType"] = []string{signType}
	data["curtime"] = []string{curTime}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", apiUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return TextTranslationResp{}, fmt.Errorf("http request fail:%w", err)
	}
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	httpResp, err := client.config.Client.PostForm(apiUrl, data)
	defer httpResp.Body.Close()
	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return TextTranslationResp{}, fmt.Errorf("get json fail,%w", err)
	}
	if httpResp.StatusCode != http.StatusOK {
		return TextTranslationResp{}, fmt.Errorf("bad result:%v", string(body))
	}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return TextTranslationResp{}, fmt.Errorf("unmarshal json fail,%w", err)
	}
	return resp, nil
}

//NewClient 初始化
func NewClient(config Config) Client {
	return Client{config: config}
}
