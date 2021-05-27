package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/spf13/cobra"
)

const (
	signType = "v3"
	APIURI   = "https://openapi.youdao.com/api"
)

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
type youDaoClient struct {
	youDaoApi string
	config    *Config
	data      url.Values
}

type Context struct {
	fromLang string
	toLang   string
	q        string
}
type Client struct {
	config *Config
}

func innerMain() int {
	var fromLang string
	var toLang string
	var command = &cobra.Command{Use: "DictionaryGo [word]", Short: "translate words",
		Long: `translate words to other language by cmdline`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			q := strings.Join(args, " ")
			youDaoClient := &youDaoClient{
				APIURI,
				InitConfig(),
				make(url.Values, 0),
			}
			client := youDaoClient.newClient(youDaoClient.config)
			context := Context{fromLang, toLang, q}
			response, err := client.TextTranslation(context, youDaoClient)

			if err != nil {
				fmt.Println("")
			}
			PrintTranslation(&response, os.Stdout)
			return
		}}
	command.Flags().StringVarP(&fromLang, "from", "f", "auto", "translate from this language")
	command.Flags().StringVarP(&toLang, "to", "t", "auto", "translate to this language")
	err := command.Execute()
	if err != nil {
		log.Fatal(err)
	}
	return 1
}

func main() {
	os.Exit(innerMain())
}

//newClient 初始化
func (youDaoClient *youDaoClient) newClient(config *Config) Client {
	var client Client
	client.config = config
	return client
}

//TextTranslation 调用API，查询翻译结果
func (client *Client) TextTranslation(ctx Context, req *youDaoClient) (DictionaryRespJson, error) {
	uuidRandNum := uuid.NewV4()

	input := TruncateString(ctx.q)
	stamp := time.Now().Unix()
	instr := client.config.AppKey + input + uuidRandNum.String() + strconv.FormatInt(stamp, 10) + client.config.AppSecret
	sign := sha256.Sum256([]byte(instr))
	signToStr := HexNumToString(sign[:])

	req.data = make(url.Values, 0)
	req.data["q"] = []string{ctx.q}
	req.data["from"] = []string{ctx.fromLang}
	req.data["to"] = []string{ctx.toLang}
	req.data["appKey"] = []string{client.config.AppKey}
	req.data["salt"] = []string{uuidRandNum.String()}
	req.data["sign"] = []string{signToStr}
	req.data["signType"] = []string{signType}
	req.data["curtime"] = []string{strconv.FormatInt(stamp, 10)}

	resp, err := http.PostForm(req.youDaoApi, req.data)
	if err != nil || resp == nil {
		fmt.Printf("open api failed,%v\n", err)
	}
	defer resp.Body.Close()
	var jsonContent DictionaryRespJson
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return DictionaryRespJson{}, err
	}

	err = json.Unmarshal(body, &jsonContent)
	if err != nil {
		return DictionaryRespJson{}, err
	}
	return jsonContent, nil
}

//PrintTranslation 将报文内容结构化显示
func PrintTranslation(jsonContent *DictionaryRespJson, writer io.Writer) {
	if jsonContent.ErrorCode != "0" {
		switch jsonContent.ErrorCode {
		case "102":
			_, _ = fmt.Fprintln(writer, "Language type not supported")
		case "103":
			_, _ = fmt.Fprintln(writer, "Translation text too long")
		case "108":
			_, _ = fmt.Fprintln(writer, "Invalid AppID,please confirm you ID")
		case "113":
			_, _ = fmt.Fprintln(writer, "The text to be queried is empty")
		case "401":
			_, _ = fmt.Fprintln(writer, "The account is overdue. Please recharge the account")
		default:
			_, _ = fmt.Fprintln(writer, "Unknown error")
		}
		return
	}
	_, _ = fmt.Fprintln(writer, "----", jsonContent.Query, "----")
	if jsonContent.Basic.UkPhonetic != "" {
		_, _ = fmt.Fprintln(writer, "英式发音:", "[ ", jsonContent.Basic.UkPhonetic, " ]")
	}
	if jsonContent.Basic.UsPhonetic != "" {
		_, _ = fmt.Fprintln(writer, "美式发音:", "[ ", jsonContent.Basic.UsPhonetic, " ]")
	}
	_, _ = fmt.Fprintln(writer, "[ 翻译结果 ]")
	for k, v := range jsonContent.Translation {
		_, _ = fmt.Fprintln(writer, "\t", k+1, ".", v)
	}
	if jsonContent.Basic.Explains != nil {
		_, _ = fmt.Fprintln(writer, "[ 网络释义 ]")
		for k, v := range jsonContent.Basic.Explains {
			_, _ = fmt.Fprintln(writer, "\t", k+1, ".", v)
		}
	}
	if jsonContent.Web != nil {
		_, _ = fmt.Fprintln(writer, "[ 延伸释义 ]")
		for k, v := range jsonContent.Web {
			_, _ = fmt.Fprintln(writer, "\t", k+1, ".", v.Key)
			_, _ = fmt.Fprint(writer, "\t翻译:")
			for _, val := range v.Value {
				_, _ = fmt.Fprint(writer, val, ",")
			}
			_, _ = fmt.Fprint(writer, "\n")
		}
	}
	if jsonContent.WebDict != nil {
		_, _ = fmt.Fprintln(writer, "[ 查看详情:]")
		for _, v := range jsonContent.WebDict {
			_, _ = fmt.Fprintln(writer, "\t", v)
		}
	}

}

//InitConfig 加载配置文件
func InitConfig() *Config {
	var config Config
	curFilePath := GetFilePath()
	content, err := ioutil.ReadFile(curFilePath + "/config.json")
	if err != nil {
		fmt.Println("config file open failed")
	}
	err1 := json.Unmarshal(content, &config)
	if err1 != nil {
		fmt.Println("config unmarshal failed")
	}
	return &config
}

//GetFilePath 获取当前文件地址
func GetFilePath() string {
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println("get file's path failed")
		return ""
	}
	newPath := strings.Replace(path, "\\", "/", -1)
	return newPath
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
