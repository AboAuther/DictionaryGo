package main

import (
	"fmt"
	"io"
)

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
