package main

import (
	"context"
	"encoding/json"
	"github.com/chromedp/chromedp/runner"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/chromedp/chromedp"
)

type info struct {
	Username string `json:username`
	Password string `json:password`
}

var userInfo info

func init() {
	// 初始化用户信息
	content, err := ioutil.ReadFile("./info.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(content, &userInfo)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	var err error

	// create context
	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create chrome instance
	//c, err := chromedp.New(ctxt, chromedp.WithLog(log.Printf))
	c, err := chromedp.New(ctxt, chromedp.WithRunnerOptions(
		runner.Path("C:\\Users\\zhen_wang\\AppData\\Local\\Google\\Chrome\\Application\\chrome.exe"),
	))
	if err != nil {
		log.Fatal(err)
	}

	// run task list
	var res string
	err = c.Run(ctxt, text(&res))
	if err != nil {
		log.Fatal(err)
	}

	// shutdown chrome
	err = c.Shutdown(ctxt)
	if err != nil {
		log.Fatal(err)
	}

	// wait for chrome to finish
	err = c.Wait()
	if err != nil {
		log.Fatal(err)
	}

	ioutil.WriteFile("./resault.csv", []byte(dataParse(res)), os.ModePerm)
	log.Printf("overview: %s", res)
}

func text(res *string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate("https://trello.com/login"),
		chromedp.SetValue("#user", userInfo.Username, chromedp.NodeVisible, chromedp.ByQuery),
		chromedp.SetValue("#password", userInfo.Password, chromedp.NodeVisible, chromedp.ByQuery),
		chromedp.Click("#login", chromedp.NodeVisible, chromedp.ByQuery),
		chromedp.Click(".board-tile-details-name[title='技术分享主题池']>div", chromedp.NodeVisible, chromedp.ByQuery),
		chromedp.Click(".show-more", chromedp.NodeVisible, chromedp.ByQuery),
		chromedp.Text(".phenom-list", res, chromedp.NodeVisible, chromedp.ByQuery),
	}
}

func dataParse(data string) string {
	dataArr := strings.Split(data, "发送中…")
	retArr := []string{}
	reg := regexp.MustCompile(`\s*?(.+)在\s+?(\d+.+?)(\d+月\d+日)`)
	for _, item := range dataArr {
		if reg.MatchString(item) {
			retArr = append(retArr, reg.ReplaceAllString(item, "$1,$2,$3"))
		}
	}
	return strings.Join(retArr, "\n")
}
