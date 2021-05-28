package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

//Config 配置文件
type Config struct {
	AppKey    string `json:"AppKey"`
	AppSecret string `json:"AppSecret"`
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
