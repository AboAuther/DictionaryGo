package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	youdao "DictionaryGo/youdao"
	"github.com/spf13/cobra"
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

func callYouDaoApi(fromLang, toLang, q string) (youdao.TextTranslationResp, error) {
	config := InitConfig()
	signToStr, salt, curTime := GetSign(q, config)
	client := youdao.NewClient((*youdao.Config)(config))
	textTranslateReq := youdao.NewTextTranslateReq(fromLang, toLang, q, signToStr, salt, curTime)
	response, err := client.TextTranslation(context.Background(), textTranslateReq)
	if err != nil {
		return youdao.TextTranslationResp{}, fmt.Errorf("%w", err)
	}
	return response, nil
}
