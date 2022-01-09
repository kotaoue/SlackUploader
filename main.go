package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/slack-go/slack"
)

var (
	userID   = flag.String("UserID", "", "Slack UserID")
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
		fmt.Printf("MsgID: %s, Text:%s, ReplyCount:%d\n", v.ClientMsgID, v.Text, v.ReplyCount)
	}

	// respChannel, respTimestamp, err := api.PostMessage(*chanelID, slack.MsgOptionText("Hello World", true))
	respTimestamp, err := api.PostEphemeral(*chanelID, *userID, slack.MsgOptionText("Hello World", true))
	fmt.Printf("respTimestamp:%s\n", respTimestamp)

	respTimestamp, err = api.PostEphemeral(*chanelID, *userID,
		slack.MsgOptionText("Hello World2", true),
		slack.MsgOptionTS(respTimestamp),
	)
	fmt.Printf("respTimestamp:%s\n", respTimestamp)

	return nil
}
