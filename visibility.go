package main

type Visibility string

const (
	VisibilityPrivate         Visibility = "private"
	VisibilityPublic          Visibility = "public"
	VisibilityPublicWhitelist Visibility = "public-whitelist"
	VisibilityESP             Visibility = "esp"
	VisibilityIAP             Visibility = "iap"
	VisibilityApigee          Visibility = "apigee"

	VisibilityUnknown Visibility = ""
)
