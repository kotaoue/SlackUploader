package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/slack-go/slack"
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
	api := slack.New(getToken())
	user, err := api.GetUserInfo("U023BECGF")
	if err != nil {
		return err
	}
	fmt.Printf("ID: %s, Fullname: %s, Email: %s\n", user.ID, user.Profile.RealName, user.Profile.Email)

	return nil
}

func getToken() string {
	fmt.Println("please input token")

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		if s.Text() != "" {
			break
		}
	}

	return s.Text()
}
