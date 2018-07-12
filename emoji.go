package main

import (
	"io/ioutil"
	"net/http"
	"time"
)

type Emoji struct {
	Name string
	URL  string
	Type string // png, gif
}

var client = &http.Client{
	Timeout: time.Second * 60,
}

func (e *Emoji) Download() ([]byte, error) {
	res, err := client.Get(e.URL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	bs, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	// TODO: check image type
	return bs, nil
}
