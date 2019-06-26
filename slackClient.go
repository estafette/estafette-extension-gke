package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/sethgrid/pester"
)

var (
	webhookURL = "https://hooks.slack.com/services/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"
)

func sendNotifications(status string) {
	var message = ""
	switch status {
	case "succeeded":
		message = "good"
	case "failed":
		message = "Your last deployment generate too many errors... rolling back"
	}

	//send slack notification
	var channels [1]string

	for i := range channels {
		err := sendSlackNotification(channels[i], "To many errors!", message, status)
		log.Fatal(err)
	}
}

func sendSlackNotification(channel, title, message, status string) (err error) {

	var requestBody io.Reader

	color := ""
	switch status {
	case "succeeded":
		color = "good"
	case "failed":
		color = "danger"
	}

	slackMessageBody := SlackMessageBody{
		Channel:  channel,
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
		log.Printf("Failed marshalling SlackMessageBody: %v. Error: %v", slackMessageBody, err)
		return
	}
	requestBody = bytes.NewReader(data)

	client := pester.New()
	client.MaxRetries = 3
	client.Backoff = pester.ExponentialJitterBackoff
	client.KeepLog = true
	request, err := http.NewRequest("POST", webhookURL, requestBody)
	if err != nil {
		log.Printf("Failed creating http client: %v", err)
		return
	}

	// add headers
	request.Header.Add("Content-type", "application/json")

	// perform actual request
	response, err := client.Do(request)
	if err != nil {
		log.Printf("Failed performing http request to Slack: %v", err)
		return
	}

	defer response.Body.Close()

	return
}
