package main

import (
	"github.com/BeepBoopHQ/go-slackbot"
	"os"
	"github.com/nlopes/slack"
	"context"
)

func main() {
	bot := slackbot.New(os.Getenv("SLACK_KEY"))

	toMe := bot.Messages(slackbot.DirectMessage, slackbot.DirectMention).Subrouter()
	toMe.Hear("(?)").MessageHandler(HelloHandler)
	bot.Run()
}

func HelloHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	bot.Reply(evt, "Oh hello!", true)
}