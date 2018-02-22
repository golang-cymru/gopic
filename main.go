package main

import (
	"github.com/BeepBoopHQ/go-slackbot"
	"os"
	"github.com/nlopes/slack"
	"context"
	"log"
	"net/http"
	"fmt"
	"io/ioutil"
)

func main() {
	bot := slackbot.New(os.Getenv("SLACK_KEY"))

	toMe := bot.Messages(slackbot.DirectMention, slackbot.DirectMessage).Subrouter()
	toMe.Hear("(?)").MessageHandler(HelloHandler)
	bot.Run()
}

func HelloHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	if evt.Upload {
		client := &http.Client{}
		req, err := http.NewRequest("GET", evt.File.URLPrivate, nil)
		if err != nil {
			log.Fatalln(err)
		}
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("SLACK_KEY")))
		res, err := client.Do(req)
		if err != nil{
			log.Fatal(err)
		}

		bytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatalln(err)
		}

		ioutil.WriteFile("upload.jpg", bytes, 0755)
	}
}