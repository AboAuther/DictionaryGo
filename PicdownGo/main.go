package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/cobra"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/sync/errgroup"
)

var (
	baseDir string
)
var configuration Parser

func worker(destDir, website string, linkChan chan string) (err error) {
	for picUrl := range linkChan {
		if strings.HasPrefix(picUrl, "//") {
			picUrl = "https:" + picUrl
		} else {
			picUrl = "https://" + website + picUrl
		}
		resp, err := http.Get(picUrl)
		if err != nil {
			return fmt.Errorf("Get picture URL failed,%w,picUrl:%s", err, picUrl)
		}
		defer resp.Body.Close()
		out, err := os.Create(destDir + "/" + fmt.Sprint(rand.Int()) + ".jpg")
		if err != nil {
			return fmt.Errorf("create file failed,%w", err)
		}
		defer out.Close()
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return fmt.Errorf("down picture failed,%w", err)
		}
	}
	return nil
}

func findDomainByUrl(postUrl string) (*Website, error) {
	var targetDomain string
	u, err := url.Parse(postUrl)
	if err != nil {
		return nil, fmt.Errorf("URL may be wrong,%w", err)
	}
	targetDomain = u.Host
	for index, website := range configuration.SupportWebsites {
		if strings.Contains(targetDomain, website.Website) {
			fmt.Println("use ", configuration.SupportWebsites[index].Website, "parser")
			return &configuration.SupportWebsites[index], nil
		}
	}
	fmt.Println("Url input is not in support list")
	return nil, nil
}

func crawler(postUrl string, workNum int) (err error) {
	res, err := http.Get(postUrl)
	if err != nil {
		return fmt.Errorf("url input is unvalued,%w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("Status code error:%d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return fmt.Errorf("create new documents failed,%w", err)
	}
	targetSiteConf, err := findDomainByUrl(postUrl)
	if err != nil {
		return fmt.Errorf("find domian failed,%w", err)
	}
	if targetSiteConf == nil {
		return nil
	}
	title := doc.Find(targetSiteConf.TitlePattern).Text()
	fmt.Println("[", targetSiteConf.Website, "]:", title, "start downloading...")
	dir := fmt.Sprintf("%v/%v - %v", baseDir, targetSiteConf.Website, title)
	_ = os.MkdirAll(dir, 0755)
	linkChan := make(chan string)
	var g errgroup.Group
	for i := 0; i < workNum; i++ {
		g.Go(func() error {
			err := worker(dir, targetSiteConf.Website, linkChan)
			return err
		})
	}
	if err := g.Wait(); err != nil {
		return fmt.Errorf(" fetch URLs failed.%w", err)
	}
	return nil
}

func checkFileIsExits(jsonFileAddr string) bool {
	if _, err := os.Stat(jsonFileAddr); err == nil {
		return true
	}
	return false
}
func reloadParser(jsonFileAddr string) (err error) {
	var file []byte
	if checkFileIsExits(jsonFileAddr) {
		file, _ = os.ReadFile(jsonFileAddr)
		fmt.Println("Loading local json file...")
	} else {
		err := errors.New("Json file is not Existed.")
		return err
	}
	if len(file) == 0 {
		fmt.Println("Json file is empty,Loading default parser...")
		file = DefaultJson
	}
	err = json.Unmarshal(file, &configuration)
	if err != nil {
		return fmt.Errorf("json file unmarshal failed,%w", err)
	}
	return nil
}
func main() {
	baseDir = "Pic"
	err := os.MkdirAll(baseDir, 0755)
	if err != nil {
		_ = fmt.Errorf("create dictionary failed,%w", err)
	}
	jsonFileAddr := fmt.Sprintf("./parser.json")
	err = reloadParser(jsonFileAddr)
	if err != nil {
		_ = fmt.Errorf("Loading jsonfile failed,%w", err)
	}
	var workNum int
	var postUrl string
	var command = &cobra.Command{
		Use:   "PicdownGo -u URL",
		Short: "Download pictures",
		Long:  `Download given url pictures by cmdline`,
		//Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if postUrl == "" {
				fmt.Println("Please use 'PicdownGo -u URL'.")
			}
			err = crawler(postUrl, workNum)
			if err != nil {
				return fmt.Errorf("%w", err)
			}
			return nil
		}}
	command.Flags().StringVarP(&postUrl, "URL", "u", "", "URL of post")
	command.Flags().IntVarP(&workNum, "workerNum", "w", 20, "numben of workers")
	err = command.Execute()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

}
