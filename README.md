# SlackUploader
1. This script uploads jpg or png to Slack: Recursively from the current directory.
1. Script is delete messages after upload.
1. And echo Permalinks.

If you want to create a [Slackbot](https://slack.com/intl/en-in/help/articles/202026038-An-introduction-to-Slackbot) that displays images randomly: Try this script.

## Requirement
* [slack-go/slack](https://github.com/slack-go/slack)
* [slack api:Access tokens](https://api.slack.com/authentication/token-types)
* [slack api:Permission scopes](https://api.slack.com/scopes)
  * [files:write](https://api.slack.com/scopes/files:write)

## Preparation
```
curl -XPOST "https://slack.com/api/auth.test?token=<TOKEN>&pretty=1"
```

## Usage
```
go run main.go -UserID=<Slack UserID> -ChannelID=<Upload target ChannelID> -token=<Bot User OAuth Token>
```