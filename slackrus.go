// Package slackrus provides a Hipchat hook for the logrus loggin package.
package slackrus

import (
	"github.com/Sirupsen/logrus"
	"github.com/johntdyer/slack-go"
)

// Project version
const (
	VERISON = "0.0.2"
)

var (
	client *slack.Client
)

// SlackrusHook is a logrus Hook for dispatching messages to the specified
// channel on Slack.
type SlackrusHook struct {
	// Messages with a log level not contained in this array
	// will not be dispatched. If nil, all messages will be dispatched.
	AcceptedLevels []logrus.Level
	HookURL        string
	IconURL        string
	Channel        string
	IconEmoji      string
	Username       string
	c              *slack.Client
}

// Levels sets which levels to sent to slack
func (sh *SlackrusHook) Levels() []logrus.Level {
	if sh.AcceptedLevels == nil {
		return AllLevels
	}
	return sh.AcceptedLevels
}

// Fire -  Sent event to slack
func (sh *SlackrusHook) Fire(e *logrus.Entry) error {
	if sh.c == nil {
		if err := sh.initClient(); err != nil {
			return err
		}
	}

	color := ""
	switch e.Level {
	case logrus.DebugLevel:
		color = "#9B30FF"
	case logrus.InfoLevel:
		color = "good"
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		color = "danger"
	default:
		color = "warning"
	}

	msg := &slack.Message{
		Username: sh.Username,
		Channel:  sh.Channel,
	}

	msg.IconEmoji = sh.IconEmoji
	msg.IconUrl = sh.IconURL

	attach := msg.NewAttachment()

	// If there are fields we need to render them at attachments
	if len(e.Data) > 0 {

		// Add a header above field data
		attach.Text = "Message fields"

		for k, v := range e.Data {
			str := toString(v)
			if str != "" {
				slackField := &slack.Field{
					Title: k,
					Value: str,
					Short: len(str) <= 20,
				}

				attach.AddField(slackField)
			}
		}
		attach.Pretext = e.Message
	} else {
		attach.Text = e.Message
	}
	attach.Fallback = e.Message
	attach.Color = color

	return sh.c.SendMessage(msg)

}

func (sh *SlackrusHook) initClient() error {
	sh.c = &slack.Client{sh.HookURL}

	if sh.Username == "" {
		sh.Username = "SlackRus"
	}

	return nil
}
