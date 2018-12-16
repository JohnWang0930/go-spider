package main

import (
	"context"
	"log"

	"github.com/chromedp/chromedp"
)

func main() {
	var err error

	// create context
	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create chrome instance
	//c, err := chromedp.New(ctxt, chromedp.WithLog(log.Printf))
	c, err := chromedp.New(ctxt)
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

	log.Printf("overview: %s", res)
}

func text(res *string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate("https://trello.com/login"),
		chromedp.SetValue("#user", "", chromedp.NodeVisible, chromedp.ByQuery),
		chromedp.SetValue("#password", "", chromedp.NodeVisible, chromedp.ByQuery),
		chromedp.Click("#login", chromedp.NodeVisible, chromedp.ByQuery),
		chromedp.Click(".board-tile-details-name[title='技术分享主题池']>div", chromedp.NodeVisible, chromedp.ByQuery),
		chromedp.Text(".js-menu-action-list", res, chromedp.NodeVisible, chromedp.ByQuery),
	}
}
