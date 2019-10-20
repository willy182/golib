package golib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// Attachment model
type Attachment struct {
	Attachments []Payload `json:"attachments"`
}

// Payload model
type Payload struct {
	Text   string  `json:"text"`
	Color  string  `json:"color"`
	Fields []Field `json:"fields"`
}

// Field model
type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

const (
	successColor = "#36a64f"
	errorColor   = "#f44b42"
)

func getCaller() string {
	var name, file string
	var line int
	var pc [16]uintptr

	n := runtime.Callers(4, pc[:])
	for _, pc := range pc[:n] {
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			continue
		}
		file, line = fn.FileLine(pc)
		name = fn.Name()
		if !strings.HasPrefix(name, "runtime.") {
			break
		}
	}

	var source string
	switch {
	case name != "":
		source = fmt.Sprintf("%v:%v", name, line)
	case file != "":
		source = fmt.Sprintf("%v:%v", file, line)
	default:
		source = fmt.Sprintf("pc:%x", pc)
	}
	return source
}

// SendNotification to slack channel
func SendNotification(title, body, ctx string, err error) {
	isActive, _ := strconv.ParseBool(os.Getenv("SLACK_NOTIFIER"))
	if !isActive {
		return
	}

	stackTrace := getCaller()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				LogError(fmt.Errorf("%v", r), "send_notification_slack", title)
			}
		}()
		message := fmt.Sprintf("*%s*\n\n%s", title, body)

		var slackPayload Payload
		slackPayload.Text = message
		slackPayload.Color = successColor
		if err != nil {
			slackPayload.Color = errorColor
			slackPayload.Text = fmt.Sprintf("%s\n*Error*: ```%s```", message, err.Error())
		}

		hostName, _ := os.Hostname()
		now := time.Now().Format(time.RFC3339)
		slackPayload.Fields = []Field{
			Field{
				Title: "Server",
				Value: hostName,
				Short: true,
			},
			Field{
				Title: "Environment",
				Value: os.Getenv("SERVER_ENV"),
				Short: true,
			},
			Field{
				Title: "Context",
				Value: ctx,
				Short: true,
			},
			Field{
				Title: "Time",
				Value: now,
				Short: true,
			},
		}

		if err != nil {
			slackPayload.Fields = append(slackPayload.Fields, Field{
				Title: "Error Line Stack",
				Value: fmt.Sprintf("`%s`", stackTrace),
				Short: true,
			})
		}

		var slackAttachment Attachment
		slackAttachment.Attachments = append(slackAttachment.Attachments, slackPayload)

		buffer := &bytes.Buffer{}
		encoder := json.NewEncoder(buffer)
		encoder.SetEscapeHTML(true)
		encoder.Encode(slackAttachment)

		url := os.Getenv("SLACK_URL")
		req, _ := http.NewRequest("POST", url, buffer)
		defer req.Body.Close()

		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}
		Log(ErrorLevel, string(body), "slack.SendNotification", "send_to_slack")
	}()
}
