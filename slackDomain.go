package main

// SlackMessageBody represents the body to send to slack
type SlackMessageBody struct {
	Channel     string                   `json:"channel"`
	Username    string                   `json:"username"`
	Text        string                   `json:"text,omitempty"`
	Attachments []SlackMessageAttachment `json:"attachments,omitempty"`
}

// SlackMessageAttachment represents attachments to a slack message
type SlackMessageAttachment struct {
	Fallback   string               `json:"fallback"`
	Color      string               `json:"color,omitempty"`
	Pretext    string               `json:"pretext,omitempty"`
	AuthorName string               `json:"author_name,omitempty"`
	AuthorLink string               `json:"author_link,omitempty"`
	AuthorIcon string               `json:"author_icon,omitempty"`
	Title      string               `json:"title,omitempty"`
	TitleLink  string               `json:"title_link,omitempty"`
	Text       string               `json:"text,omitempty"`
	ImageURL   string               `json:"image_url,omitempty"`
	ThumbURL   string               `json:"thumb_url,omitempty"`
	Footer     string               `json:"footer,omitempty"`
	FooterIcon string               `json:"footer_icon,omitempty"`
	Ts         int                  `json:"ts,omitempty"`
	MarkdownIn []string             `json:"mrkdwn_in,omitempty"`
	Actions    []SlackMessageAction `json:"actions,omitempty"`
}

// SlackMessageAction represents an action (button)
type SlackMessageAction struct {
	Type    string `json:"type,omitempty"`
	Name    string `json:"name,omitempty"`
	Text    string `json:"text,omitempty"`
	URL     string `json:"url,omitempty"`
	Style   string `json:"style,omitempty"`
	Confirm string `json:"confirm,omitempty"`
}