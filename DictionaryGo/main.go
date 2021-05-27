package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func innerMain() int {
	var fromLang string
	var toLang string
	var command = &cobra.Command{Use: "DictionaryGo [word]", Short: "translate words",
		Long: `translate words to other language by cmdline`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			q := strings.Join(args, " ")
			response, err := callYouDaoApi(fromLang, toLang, q)
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

func callYouDaoApi(fromLang, toLang, q string) (DictionaryRespJson, error) {
	youDaoClient := &youDaoClient{
		APIURI,
		InitConfig(),
		make(url.Values, 0),
	}
	client := youDaoClient.newClient(youDaoClient.config)
	context := Context{fromLang, toLang, q}
	response, err := client.TextTranslation(context, youDaoClient)
	if err != nil {
		return DictionaryRespJson{}, fmt.Errorf("%w", err)
	}
	return response, nil
}

//newClient 初始化
func (youDaoClient *youDaoClient) newClient(config *Config) Client {
	var client Client
	client.config = config
	return client
}

//TextTranslation 调用API，查询翻译结果
func (client *Client) TextTranslation(ctx Context, req *youDaoClient) (DictionaryRespJson, error) {

	signToStr, salt, curTime := GetSign(ctx.q, client.config)
	req.data = make(url.Values, 0)
	req.data["q"] = []string{ctx.q}
	req.data["from"] = []string{ctx.fromLang}
	req.data["to"] = []string{ctx.toLang}
	req.data["appKey"] = []string{client.config.AppKey}
	req.data["salt"] = []string{salt}
	req.data["sign"] = []string{signToStr}
	req.data["signType"] = []string{signType}
	req.data["curtime"] = []string{curTime}

	resp, err := http.PostForm(req.youDaoApi, req.data)
	if err != nil {
		return DictionaryRespJson{}, fmt.Errorf("%w", err)
	}
	defer resp.Body.Close()
	var jsonContent DictionaryRespJson
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return DictionaryRespJson{}, fmt.Errorf("%w", err)
	}

	err = json.Unmarshal(body, &jsonContent)
	if err != nil {
		return DictionaryRespJson{}, fmt.Errorf("%w", err)
	}
	return jsonContent, nil
}
