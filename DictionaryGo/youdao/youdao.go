package youdao

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
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
	FromLang string
	ToLang   string
	Q        string
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
	signToStr, salt, curTime := GetSign(req.Q, client.config)
	data := make(url.Values, 0)
	data["q"] = []string{req.Q}
	data["from"] = []string{req.FromLang}
	data["to"] = []string{req.ToLang}
	data["appKey"] = []string{client.config.AppKey}
	data["salt"] = []string{salt}
	data["sign"] = []string{signToStr}
	data["signType"] = []string{signType}
	data["curtime"] = []string{curTime}

	cli := http.Client{Timeout: 5 * time.Second}
	respon, err := cli.PostForm(ApiUrl, data)
	if err != nil {
		return TextTranslationResp{}, fmt.Errorf("http request fail%w", err)
	}
	defer respon.Body.Close()
	body, err := ioutil.ReadAll(respon.Body)
	if err != nil {
		return TextTranslationResp{}, fmt.Errorf("get json fail,%w", err)
	}

	err = json.Unmarshal(body, &resp)
	if err != nil {
		return TextTranslationResp{}, fmt.Errorf("unmarshal json fail,%w", err)
	}
	return resp, nil
}

//newClient 初始化
func NewClient(config *Config) Client {
	var client Client
	client.config = config
	return client
}
func NewTextTranslateReq(fromLang, toLang, q string) (ttr TextTranslationReq) {
	ttr.ToLang = toLang
	ttr.FromLang = fromLang
	ttr.Q = q
	return
}
func GetSign(q string, config *Config) (signToStr, salt, curTime string) {
	uuidRandNum := uuid.NewV4()
	salt = uuidRandNum.String()
	input := TruncateString(q)
	stamp := time.Now().Unix()
	curTime = strconv.FormatInt(stamp, 10)
	instr := config.AppKey + input + salt + strconv.FormatInt(stamp, 10) + config.AppSecret
	sign := sha256.Sum256([]byte(instr))
	signToStr = HexNumToString(sign[:])
	return
}

//TruncateString 截断查询字符串(有道官方要求字符串长度大于20则取前10+字符串长度+后10，否则返回q)
func TruncateString(q string) string {
	res := make([]rune, 10)
	temp := []rune(q)
	qLen := len(temp)
	if qLen <= 20 {
		return q
	} else {
		res = temp[:10]
		strQLen := strconv.Itoa(qLen)
		res = append(res, []rune(strQLen)...)
		res = append(res, temp[qLen-10:qLen]...)
		return string(res)
	}
}

//HexNumToString 将十六进制Sign转换为字符串
func HexNumToString(hexnum []byte) (res string) {
	for _, v := range hexnum {
		str := strconv.FormatUint(uint64(v), 16)
		if len(str) == 1 {
			res = res + "0" + str
		} else {
			res += str
		}
	}
	return res
}
