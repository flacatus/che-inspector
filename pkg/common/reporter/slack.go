package reporter

import (
	"fmt"
	"github.com/flacatus/che-inspector/pkg/api"
	"github.com/slack-go/slack"
)

var slackChannelCache = make(map[string]slack.Channel)

// SendSlackMessage sends the message text to the provided channel (if it can be found).
func SendSlackMessage(reporter *api.CheReporterSpec, message string) error {
	token := reporter.Token
	if len(token) == 0 {
		return fmt.Errorf("no slack token configured")
	}
	slackAPI := slack.New(token)
	var slackChannel slack.Channel
	var ok bool

	if slackChannel, ok = slackChannelCache[reporter.Channel]; !ok {
		cursor := ""
		for {
			channels, nextCursor, err := slackAPI.GetConversations(&slack.GetConversationsParameters{Cursor: cursor})
			if err != nil {
				return err
			}
			for _, c := range channels {
				slackChannelCache[c.Name] = c
				if c.Name == reporter.Channel {
					slackChannel = c
				}
			}
			if nextCursor == "" {
				break
			}
			cursor = nextCursor
		}
	}

	if slackChannel.ID == "" {
		return fmt.Errorf("no slack channel named `%s` found", reporter.Channel)
	}

	_, _, err := slackAPI.PostMessage(slackChannel.ID, slack.MsgOptionText(message, false))
	return err
}
