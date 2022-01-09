package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/slack-go/slack"
)

var (
	userID   = flag.String("UserID", "", "Slack UserID")
	chanelID = flag.String("ChannelID", "", "upload target ChanelID cf.https://api.slack.com/methods/channels.list/test")
	token    = flag.String("token", "", "Bot User OAuth Token cf.https://api.slack.com/authentication/token-types")
	delete   = flag.Bool("delete", false, "delete latest message")
	force    = flag.Bool("force", false, "do not confirm")
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

	var permaLinks []string
	err := filepath.Walk("./", func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		switch filepath.Ext(path) {
		case ".png", ".jpg":
			if *force || confirm(path) {
				p, err := upload(api, path)
				if err != nil {
					return err
				}
				permaLinks = append(permaLinks, p)
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	for _, v := range permaLinks {
		fmt.Println(v)
	}

	return nil
}

func upload(api *slack.Client, path string) (string, error) {
	fmt.Printf("Upload \"%s\"\n", path)

	fp, err := os.Open(path)
	if err != nil {
		return "", err
	}

	file, err := api.UploadFile(
		slack.FileUploadParameters{
			Reader:   fp,
			Filename: fp.Name(),
			Channels: []string{*chanelID},
		},
	)

	_, _, err = api.DeleteMessage(*chanelID, file.Shares.Public[*chanelID][0].Ts)
	if err != nil {
		return "", err
	}

	return file.Permalink, nil
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

func confirm(path string) bool {
	fmt.Printf("Do you really want to upload \"%s\"? [Y/n]\n", path)

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		switch s.Text() {
		case "Y":
			return true
		case "n":
			return false
		}
	}

	return false
}
