package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if len(os.Args) < 3 || os.Args[1] != "upload" {
		log.Fatal("usage: ./slack-emoji upload /path/to/emoji")
	}
	err := godotenv.Load("config.env")
	if err != nil {
		log.Println(err)
	}
	uploader := &Uploader{
		URL:    fmt.Sprintf("https://%s.slack.com", os.Getenv("SLACK_WORKSPACE")),
		Cookie: os.Getenv("SLACK_COOKIE"),
	}

	err = uploader.UploadDir(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
}

func export() {
	file, err := ioutil.ReadFile("sample.html")
	if err != nil {
		panic(err)
	}
	emoji, err := ExtractEmoji(file)
	if err != nil {
		panic(err)
	}
	Download("./emoji", emoji, 3)
}
