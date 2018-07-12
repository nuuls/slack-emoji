package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Uploader struct {
	Cookie    string
	URL       string
	csrfToken string
}

func (u *Uploader) UploadEmoji(path string) error {
	if u.csrfToken == "" {
		csrf, err := u.getCsrfToken()
		if err != nil {
			return err
		}
		u.csrfToken = csrf
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}

	buf := &bytes.Buffer{}
	mp := multipart.NewWriter(buf)
	basePath := filepath.Base(path)
	name := strings.Split(basePath, ".")[0]

	mp.WriteField("add", "1")
	mp.WriteField("crumb", u.csrfToken)
	mp.WriteField("name", name)
	mp.WriteField("mode", "data")

	fw, err := mp.CreateFormFile("img", basePath)
	if err != nil {
		return err
	}
	_, err = io.Copy(fw, file)
	if err != nil {
		return err
	}

	mp.Close()
	file.Close()

	req, err := http.NewRequest(http.MethodPost, u.URL+"/customize/emoji", buf)
	if err != nil {
		return err
	}

	req.Header.Set("Cookie", u.Cookie)
	req.Header.Set("content-type", "multipart/form-data")
	req.Header.Add("content-type", "boundary="+mp.Boundary())

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	out, _ := ioutil.ReadAll(res.Body)
	if !bytes.Contains(out, []byte("Your new emoji has been saved!")) {
		return errors.New("something broke")
	}

	return nil
}

var csrfRe = regexp.MustCompile(`<input .+ name="crumb" value="([^"]+)`)

func (u *Uploader) getCsrfToken() (string, error) {
	req, err := http.NewRequest(http.MethodGet, u.URL+"/customize/emoji", nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Cookie", u.Cookie)

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	out, _ := ioutil.ReadAll(res.Body)
	match := csrfRe.FindSubmatch(out)
	if len(match) > 0 {
		return string(match[1]), nil
	}
	return "", errors.New("no token found")
}

func (u *Uploader) UploadDir(dir string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, info := range files {
		err := u.UploadEmoji(filepath.Join(dir, info.Name()))
		if err != nil {
			log.Println("cannot upload", info.Name(), err)
		} else {
			log.Println("uploaded", info.Name())
		}
	}
	return nil
}
