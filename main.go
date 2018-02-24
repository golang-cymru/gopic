package main

import (
	"github.com/BeepBoopHQ/go-slackbot"
	"os"
	"github.com/nlopes/slack"
	"context"
	"log"
	"net/http"
	"fmt"
	"image/png"
	"image/jpeg"
	"image/gif"
	"image"
	"golang.org/x/image/bmp"
	"bytes"
	"gocv.io/x/gocv"
	"io/ioutil"
	"image/color"
	"path"
)

var classifier gocv.CascadeClassifier

func main() {
	bot := slackbot.New(os.Getenv("SLACK_KEY"))

	classifier = gocv.NewCascadeClassifier()
	defer classifier.Close()

	classifier.Load("haarcascade_frontalface_default.xml")

	toMe := bot.Messages(slackbot.DirectMention, slackbot.DirectMessage).Subrouter()
	toMe.Hear("(?)").MessageHandler(HelloHandler)
	bot.Run()
}

func HelloHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	if evt.Upload {
		client := &http.Client{}
		req, err := http.NewRequest("GET", evt.File.URLPrivate, nil)
		if err != nil {
			log.Println(err)
			return
		}
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("SLACK_KEY")))
		res, err := client.Do(req)
		if err != nil {
			log.Println(err)
			return
		}

		buf, _ := ioutil.ReadAll(res.Body)
		rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
		rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))

		var src image.Image
		_, format, err := image.DecodeConfig(rdr1)
		if err != nil {
			log.Fatalln(err)
		}
		switch format {
		case "png":
			src, err = png.Decode(rdr2)
		case "jpeg":
			src, err = jpeg.Decode(rdr2)
		case "gif":
			src, err = gif.Decode(rdr2)
		}

		if err != nil {
			log.Println(err)
			return
		}

		bmpIn := new(bytes.Buffer)
		err = bmp.Encode(bmpIn, src)
		if err != nil {
			log.Println(err)
			return
		}

		cvImg := gocv.IMDecode(bmpIn.Bytes(), gocv.IMReadColor)
		defer cvImg.Close()

		if cvImg.Empty() {
			return
		}

		rects := classifier.DetectMultiScale(cvImg)
		blue := color.RGBA{R: 0, G: 0, B: 255, A: 0}

		for _, r := range rects {
			gocv.Rectangle(cvImg, r, blue, 3)
		}

		outBmp, err := gocv.IMEncode(".bmp", cvImg)
		if err != nil {
			log.Println(err)
			return
		}
		outBmpBuf := bytes.NewBuffer(outBmp)

		outImg, err := bmp.Decode(outBmpBuf)
		if err != nil {
			log.Println(err)
			return
		}

		jpegOut := new(bytes.Buffer)
		err = jpeg.Encode(jpegOut, outImg, &jpeg.Options{Quality:95})
		if err != nil {
			log.Println(err)
			return
		}

		fileOptions := slack.FileUploadParameters{
			Filename: path.Base(evt.File.URLPrivate),
			Reader: jpegOut,
			Channels: []string{evt.Channel},
		}
		_, err = bot.Client.UploadFile(fileOptions)
		if err != nil {
			log.Println(err)
			return
		}
	}
}