package main

import (
	"bytes"
	"encoding/json"
	"github.com/AboAuther/DictionaryGo/youdao"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPrintTranslationSuccess(t *testing.T) {
	buf, err := ioutil.ReadFile("testData/work.json")
	require.Nil(t, err)
	var req youdao.TextTranslationResp
	err = json.Unmarshal(buf, &req)
	require.Nil(t, err)
	b := bytes.Buffer{}
	printTranslation(&req, &b)
	require.Equal(t, `---- work ----
英式发音: [  wɜːk  ]
美式发音: [  wɜːrk  ]
[ 翻译结果 ]
	 1 . 工作
[ 网络释义 ]
	 1 . n. 工作；功；产品；操作；职业；行为；事业；工厂；著作；文学、音乐或艺术作品
	 2 . vt. 使工作；操作；经营；使缓慢前进
	 3 . vi. 工作；运作；起作用
	 4 . n. （Work）（英、埃塞、丹、冰、美）沃克（人名）
[ 延伸释义 ]
	 1 . Work
	翻译:作品,起作用,工件,运转,
	 2 . work permit
	翻译:工作许可,工作证,工作准证,
	 3 . at work
	翻译:在工作,忙于,上班,在上班,
[ 查看详情:]
	 http://mobile.youdao.com/dict?le=eng&q=work
`, b.String())
}
func TestPrintTranslationErrCode(t *testing.T) {
	var req youdao.TextTranslationResp
	b := bytes.Buffer{}
	printTranslation(&req, &b)
	require.Equal(t, "Unknown error\n", b.String())
}
