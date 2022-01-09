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
		return deleteLatestMessage(api)
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
	_, _, err = api.DeleteMessage(*chanelID, file.Shares.Public[*chanelID][0].Ts)
	if err != nil {
		return err
	}

	return nil
}

func deleteLatestMessage(api *slack.Client) error {
	resp, err := api.GetConversationHistory(&slack.GetConversationHistoryParameters{ChannelID: *chanelID, Limit: 1})
	if err != nil {
		return err
	}

	for _, m := range resp.Messages {
		if m.ReplyCount > 0 {
			fmt.Println("has replies")
			if err := deleteReplies(api, m.Timestamp, ""); err != nil {
				return err
			}
		}

		_, _, err := api.DeleteMessage(*chanelID, m.Timestamp)
		if err != nil {
			return err
		}
	}

	return nil
}

func deleteReplies(api *slack.Client, ts, cursor string) error {
	p := &slack.GetConversationRepliesParameters{ChannelID: *chanelID, Timestamp: ts}
	if cursor != "" {
		p.Cursor = cursor
	}

	msgs, hasMore, nextCursor, err := api.GetConversationReplies(p)
	if err != nil {
		return err
	}

	for k, m := range msgs {
		// index 0 is parent message
		if k > 0 {
			fmt.Printf("MsgID: %s, Text:%s, ReplyCount:%d, TimeStamp:%s\n", m.ClientMsgID, m.Text, m.ReplyCount, m.Timestamp)
			_, _, err := api.DeleteMessage(*chanelID, m.Timestamp)
			if err != nil {
				return err
			}
		}
	}

	if hasMore {
		deleteReplies(api, ts, nextCursor)
	}

	return nil
}
