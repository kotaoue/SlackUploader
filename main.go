package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/slack-go/slack"
)

var (
	chanelID = flag.String("ChannelID", "", "Upload target ChanelID cf.https://api.slack.com/methods/channels.list/test")
	token    = flag.String("token", "", "Bot User OAuth Token cf.https://api.slack.com/authentication/token-types")
)

func init() {
	flag.Parse()
}

func main() {
	if err := Main(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func Main() error {
	api := slack.New(*token)

	resp, err := api.GetConversationHistory(&slack.GetConversationHistoryParameters{ChannelID: *chanelID, Limit: 3})
	if err != nil {
		fmt.Printf("%s\n", err)
		return nil
	}

	for _, v := range resp.Messages {
		fmt.Printf("Text:%s, Replies:%d\n", v.Text, len(v.Replies))
	}

	return nil
}
