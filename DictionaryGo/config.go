package main

import (
	"encoding/json"
	"fmt"
	"github.com/AboAuther/DictionaryGo/youdao"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

//initConfig 加载配置文件
func initConfig() (youdao.Config, error) {
	var config youdao.Config
	curFilePath := getFilePath()
	content, err := ioutil.ReadFile(curFilePath + "/config.json")
	if err != nil {
		return config, fmt.Errorf("%w", err)
	}
	err = json.Unmarshal(content, &config)
	if err != nil {
		return config, fmt.Errorf("%w", err)
	}
	return config, err
}

//getFilePath 获取当前执行二进制文件绝对目录
func getFilePath() string {
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println("get file's path failed")
		return ""
	}
	newPath := strings.Replace(path, "\\", "/", -1)
	return newPath
}
