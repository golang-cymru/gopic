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
	"path"
	"image/png"
	"image/jpeg"
	"strings"
	"image"
	"golang.org/x/image/bmp"
	"bufio"
	"bytes"
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

		ext := path.Ext(evt.File.URLPrivate)
		var src image.Image
		switch strings.ToLower(ext) {
		case "png":
			src, err = png.Decode(res.Body)
		case "jpg", "jpeg":
			src, err = jpeg.Decode(res.Body)

		if err != nil {
			log.Fatalln(err)
		}

		var bmpBuffer bytes.Buffer
		bmpOut := bufio.NewWriter(&bmpBuffer)

		err = bmp.Encode(bmpOut, src)

		ioutil.WriteFile("upload.bmp", bmpBuffer.Bytes(), 0755)
	}
}