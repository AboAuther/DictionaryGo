package main

type Parser struct {
	SupportWebsites []Website `json:"support_websites"`
}
type Website struct {
	Website        string `json:"website"`
	TitlePattern   string `json:"title_pattern"`
	ImgPattern     string `json:"img_pattern"`
	ImgAddrPattern string `json:"img_addr_pattern"`
}

var DefaultJson = []byte(`{
	"support_websites":[
		{
			"website": "cnblogs.com",
			"title_pattern": "title",
			"img_pattern" : ".navbar-branding img",
			"img_attr_pattern": "src"
		},
		{
			"website": "qq.com",
			"title_pattern": "title",
			"img_pattern" : ".pic img",
			"img_attr_pattern": "src"
		},
		{
			"website": "haicoder.net",
			"title_pattern": "title",
			"img_pattern" : ".logo img",
			"img_attr_pattern": "src"
		},
		{
			"website": "juejin.cn/post/6979532761954533390/",
			"title_pattern": "title",
			"img_pattern" : ".mobile img",
			"img_attr_pattern": "src"
		}	
	]
}`)
