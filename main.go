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
	delete   = flag.Bool("delete", false, "Delete latest message")
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

	if *delete {
		deleteLatestMessage(api)
	}

	fp, err := os.Open("icon.png")
	if err != nil {
		return err
	}

	file, err := api.UploadFile(
		slack.FileUploadParameters{
			Reader:   fp,
			Filename: fp.Name(),
			Channels: []string{*chanelID},
		},
	)

	fmt.Printf("%s\n", file.Permalink)
	api.DeleteMessage(*chanelID, file.Shares.Public[*chanelID][0].Ts)

	return nil
}

func deleteLatestMessage(api *slack.Client) error {
	resp, err := api.GetConversationHistory(&slack.GetConversationHistoryParameters{ChannelID: *chanelID, Limit: 1})
	if err != nil {
		return err
	}

	for _, v := range resp.Messages {
		fmt.Printf("MsgID: %s, Text:%s, ReplyCount:%d, TimeStamp:%s\n", v.ClientMsgID, v.Text, v.ReplyCount, v.Timestamp)
		api.DeleteMessage(*chanelID, v.Timestamp)
	}

	return nil
}
