package main

import (
	"encoding/json"
	"errors"
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
	curFilePath, err := getFilePath()
	if err != nil {
		return config, fmt.Errorf("getFilePath failed, err:%w", err)
	}
	content, err := ioutil.ReadFile(curFilePath + "/config.json")
	if err != nil {
		return config, fmt.Errorf("read config failed, err:%w", err)
	}
	err = json.Unmarshal(content, &config)
	if err != nil {
		return config, fmt.Errorf("parse config failed, err:%w", err)
	}
	return config, err
}

//getFilePath 获取当前执行二进制文件绝对目录
func getFilePath() (string, error) {
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", errors.New("get file's path failed")
	}
	newPath := strings.Replace(path, "\\", "/", -1)
	return newPath, nil
}
