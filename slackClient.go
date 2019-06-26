package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/sethgrid/pester"
)

var (
	webhookURL = "https://hooks.slack.com/services/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"
)

func sendNotifications(status string, stage string, params Params) {
	var message = ""
	var title = ""
	switch status {
	case "succeeded":
		message = "Successful deployment of " + params.App + " " + stage
		title = params.App + " in " + stage + " deployed"
	case "failed":
		message = "Your last deployment of " + params.App + " in " + stage + " generated too many errors... rolling back"
		title = "Too many errors!"
	}

	err := sendSlackNotification(title, message, status)
	logInfo(fmt.Sprintln(err))
}

func sendSlackNotification(title, message, status string) (err error) {

	var requestBody io.Reader

	color := ""
	switch status {
	case "succeeded":
		color = "good"
	case "failed":
		color = "danger"
	}

	slackMessageBody := SlackMessageBody{
		Username: "Mary Poppins",
		Attachments: []SlackMessageAttachment{
			SlackMessageAttachment{
				Fallback:   message,
				Title:      title,
				Text:       message,
				Color:      color,
				MarkdownIn: []string{"text"},
			},
		},
	}

	data, err := json.Marshal(slackMessageBody)
	if err != nil {
		logInfo("Failed marshalling SlackMessageBody: %v. Error: %v", slackMessageBody, err)
		return
	}
	requestBody = bytes.NewReader(data)

	client := pester.New()
	client.MaxRetries = 3
	client.Backoff = pester.ExponentialJitterBackoff
	client.KeepLog = true
	request, err := http.NewRequest("POST", webhookURL, requestBody)
	if err != nil {
		logInfo("Failed creating http client: %v", err)
		return
	}

	// add headers
	request.Header.Add("Content-type", "application/json")

	// perform actual request
	response, err := client.Do(request)
	if err != nil {
		logInfo("Failed performing http request to Slack: %v", err)
		return
	}

	defer response.Body.Close()

	return
}
