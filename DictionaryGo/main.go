package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/AboAuther/DictionaryGo/youdao"
	"github.com/spf13/cobra"
)

const apiTimeout = time.Second * 5

func main() {
	var fromLang string
	var toLang string
	var command = &cobra.Command{Use: "DictionaryGo [word]", Short: "translate words",
		Long: `translate words to other language by cmdline`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			q := strings.Join(args, " ")
			config, err := initConfig()
			if err != nil {
				return fmt.Errorf("can not open config file, %w", err)
			}
			config.Client = &http.Client{Timeout: apiTimeout}
			client := youdao.NewClient(config)
			response, err := client.TextTranslation(context.Background(), youdao.TextTranslationReq{FromLang: fromLang, ToLang: toLang, Q: q})
			if err != nil {
				return fmt.Errorf("call youdao api failed,%w", err)
			}
			PrintTranslation(&response, os.Stdout)
			return nil
		}}
	command.Flags().StringVarP(&fromLang, "from", "f", "en", "translate from this language")
	command.Flags().StringVarP(&toLang, "to", "t", "zh-CHS", "translate to this language")
	err := command.Execute()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}
