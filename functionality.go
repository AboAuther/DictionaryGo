package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

//Config 配置文件
type Config struct {
	AppKey    string `json:"AppKey"`
	AppSecret string `json:"AppSecret"`
}

//InitConfig 加载配置文件
func InitConfig() (*Config, error) {
	var config Config
	curFilePath := GetFilePath()
	content, err := ioutil.ReadFile(curFilePath + "/config.json")
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	err = json.Unmarshal(content, &config)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return &config, err
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
