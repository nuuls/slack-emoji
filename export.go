package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sync"
)

var trRe = regexp.MustCompile(`(?Um)<tr class="emoji_row">(?:\s|.)+</tr>`)
var emojiRe = regexp.MustCompile(`(?Um)data-original="([^"]+)"(?:\s|.)+data-emoji-name="(.+)"`)

func ExtractEmoji(page []byte) ([]Emoji, error) {
	out := []Emoji{}
	rawEmoji := trRe.FindAll(page, -1)
	for _, raw := range rawEmoji {
		m := emojiRe.FindSubmatch(raw)
		if len(m) < 1 {
			return nil, fmt.Errorf("invalid emoji %s", string(raw))
		}
		url, name := m[1], m[2]
		out = append(out, Emoji{
			URL:  string(url),
			Name: string(name),
		})
	}
	return out, nil
}

func Download(dir string, emoji []Emoji, workers int) error {
	err := os.MkdirAll(dir, 0754)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	wg.Add(len(emoji))
	for i := 0; i < workers; i++ {
		go func(worker int) {
			for j := worker; j < len(emoji); j += workers {
				defer wg.Done()
				e := emoji[j]
				log.Println("downloading:", e.Name, e.URL)
				bs, err := e.Download()
				if err != nil {
					log.Println(err)
					return
				}
				file, err := os.Create(filepath.Join(dir, e.Name+".png"))
				if err != nil {
					panic(err)
				}
				file.Write(bs)
				file.Close()
				log.Println("done:", j, "/", len(emoji))
			}
		}(i)
	}
	wg.Wait()
	return nil
}
