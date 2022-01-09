# SlackUploader
Image upload to Slack

## Requirement
* [slack-go/slack](https://github.com/slack-go/slack)
* [slack api:Access tokens](https://api.slack.com/authentication/token-types)
* [slack api:Permission scopes](https://api.slack.com/scopes)
  * [files:write](https://api.slack.com/scopes/files:write)
  * [im:write](https://api.slack.com/scopes/im:write)

## Preparation
```
curl -XPOST "https://slack.com/api/auth.test?token=<TOKEN>&pretty=1"
```
