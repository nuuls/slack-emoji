# slack-emoji

### Easily bulk upload emojis to slack

```bash
go get -u github.com/nuuls/slack-emoji
cd $GOPATH/src/github.com/nuuls/slack-emoji
go build
cp config.example.env config.env
vim config.env
./slack-emoji upload /path/to/emojis
```

To get the cookie go to https://yourworkspace.slack.com/customize/emoji , open the network tab, click on the first request and copy the Request headers `cookie` value
