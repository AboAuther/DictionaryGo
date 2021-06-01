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
				return fmt.Errorf("call youdao api fail,%w", err)
			}
			PrintTranslation(&response, os.Stdout)
			return
		}}
	command.Flags().StringVarP(&fromLang, "from", "f", "auto", "translate from this language")
	command.Flags().StringVarP(&toLang, "to", "t", "auto", "translate to this language")
	err := command.Execute()
	if err != nil {
		log.Fatal(err)
		return 1
	}
	return 0
}

func main() {
	os.Exit(innerMain())
}

func callYouDaoApi(fromLang, toLang, q string) (youdao.TextTranslationResp, error) {
	config, err := InitConfig()
	if err != nil {
		return youdao.TextTranslationResp{}, fmt.Errorf("can not open config file, %w", err)
	}
	client := youdao.NewClient((*youdao.Config)(config))
	response, err := client.TextTranslation(context.Background(), youdao.TextTranslationReq{FromLang: fromLang, ToLang: toLang, Q: q})
	if err != nil {
		return youdao.TextTranslationResp{}, fmt.Errorf("%w", err)
	}
	return response, nil
}
