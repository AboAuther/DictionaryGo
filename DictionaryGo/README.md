#DictionaryGo

###项目利用第三方API服务实现了一个利用命令行进行词句文本翻译的功能
##API使用方式
api的使用方式需要参考有道智云的官方文档
https://ai.youdao.com

##有道智云官方文档

从文档上我们获取到以下信息：

文本翻译接口地址: https://openapi.youdao.com/api

协议：

| 规则     | 描述               |
| -------- | ------------------ |
| 传输方式 | HTTPS              |
| 请求方式 | GET/POST           |
| 字符编码 | 统一使用UTF-8 编码 |
| 请求格式 | 表单               |
| 响应格式 | JSON               |

表单中的参数：

| 字段名   | 类型 | 含义                      | 必填  | 备注                                                         |
| -------- | ---- | ------------------------- | ----- | ------------------------------------------------------------ |
| q        | text | 待翻译文本                | True  | 必须是UTF-8编码                                              |
| from     | text | 源语言                    | True  | 参考下方 [支持语言](https://ai.youdao.com/DOCSIRMA/html/自然语言翻译/API文档/文本翻译服务/文本翻译服务-API文档.html#section-9) (可设置为auto) |
| to       | text | 目标语言                  | True  | 参考下方 [支持语言](https://ai.youdao.com/DOCSIRMA/html/自然语言翻译/API文档/文本翻译服务/文本翻译服务-API文档.html#section-9) (可设置为auto) |
| appKey   | text | 应用ID                    | True  | 可在 [应用管理](https://ai.youdao.com/appmgr.s) 查看         |
| salt     | text | UUID                      | True  | UUID                                                         |
| sign     | text | 签名                      | True  | sha256(应用ID+input+salt+curtime+应用密钥)                   |
| signType | text | 签名类型                  | True  | v3                                                           |
| curtime  | text | 当前UTC时间戳(秒)         | true  | TimeStamp                                                    |
| ext      | text | 翻译结果音频格式，支持mp3 | false | mp3                                                          |
| voice    | text | 翻译结果发音选择          | false | 0为女声，1为男声。默认为女声                                 |

> 签名生成方法如下：
> signType=v3；
> sign=sha256(`应用ID`+`input`+`salt`+`curtime`+`应用密钥`)；
> 其中，input的计算方式为：`input`=`q前10个字符` + `q长度` + `q后10个字符`（当q长度大于20）或 `input`=`q字符串`（当q长度小于等于20）；
## cobra

cobra是一个构建命令行工具的库，我们先大致描述一下我们需要的命令结构，首先word是必须的，还要附加两个标志(flag)：from和to。

所以大概就是这个样子：

```bash
$ ./DictionaryGo q --from en --to zh-CSH
```

或者简写成

```bash
$ ./DictionaryGo q -f en -t zh-CSH
```
默认翻译(英->中/中->英)
```bash
$ ./DictionaryGo q 
```

## 加载config.json配置

由于golang自带有json编解码库，所以我们使用json格式的配置文件。

如前面所述，配置文件需要加载appKey和appSecret两个参数，因此定义如下：

```json
{
    "appKey":"your app key",
    "appSecret":"your app secret code"
}
```

json在golang中，使用tag来指定json与结构体的映射

```go
type Config struct {
	AppKey    string `json:"appKey"`
	AppSecret string `json:"appSecret"`
}
```

json是一个文本文件，所以我们首先需要把文件中的内容读取出来,对应的函数获取内容之后使用http.PostForm发起请求，响应的结果会保存在http.Response中

```go
httpReq, err := http.NewRequestWithContext(ctx, "POST", apiUrl, strings.NewReader(data.Encode()))
if err != nil {
    return TextTranslationResp{}, fmt.Errorf("http request fail:%w", err)
}
httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
httpResp, err := client.config.Client.PostForm(apiUrl, data)
defer httpResp.Body.Close()
```

提取body中的json数据

```go
body, err := ioutil.ReadAll(httpResp.Body)
if err != nil {
    return TextTranslationResp{}, fmt.Errorf("get json fail,%w", err)
}
```

注意：此时的appKey和appSecret必须要是有效的值，否则无法得到想要的结果
## 解析查询的json结果数据

数据解析我们可以参考实际返回的json结果和有道智云的文档说明。

返回的结果是json格式，包含字段与FROM和TO的值有关，具体说明如下：

| 字段名       | 类型  | 含义             | 备注                                                         |
| ------------ | ----- | ---------------- | ------------------------------------------------------------ |
| errorCode    | text  | 错误返回码       | 一定存在                                                     |
| query        | text  | 源语言           | 查询正确时，一定存在                                         |
| translation  | Array | 翻译结果         | 查询正确时，一定存在                                         |
| basic        | text  | 词义             | 基本词典，查词时才有                                         |
| web          | Array | 词义             | 网络释义，该结果不一定存在                                   |
| l            | text  | 源语言和目标语言 | 一定存在                                                     |
| dict         | text  | 词典deeplink     | 查询语种为支持语言时，存在                                   |
| webdict      | text  | webdeeplink      | 查询语种为支持语言时，存在                                   |
| tSpeakUrl    | text  | 翻译结果发音地址 | 翻译成功一定存在，需要应用绑定语音合成实例才能正常播放 否则返回110错误码 |
| speakUrl     | text  | 源语言发音地址   | 翻译成功一定存在，需要应用绑定语音合成实例才能正常播放 否则返回110错误码 |
| returnPhrase | Array | 单词校验后的结果 | 主要校验字母大小写、单词前含符号、中文简繁体                 |

注：

a. 中文查词的basic字段只包含explains字段。

b. 英文查词的basic字段中又包含以下字段。

| 字段        | 含义                                             |
| ----------- | ------------------------------------------------ |
| us-phonetic | 美式音标，英文查词成功，一定存在                 |
| phonetic    | 默认音标，默认是英式音标，英文查词成功，一定存在 |
| uk-phonetic | 英式音标，英文查词成功，一定存在                 |
| uk-speech   | 英式发音，英文查词成功，一定存在                 |
| us-speech   | 美式发音，英文查词成功，一定存在                 |
| explains    | 基本释义                               

定义对应的字段的结构之后调用json解析并格式化打印翻译结果

###测试效果
```bash
---- work ----
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
```
