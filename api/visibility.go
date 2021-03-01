package api

type Visibility string

const (
	VisibilityPrivate         Visibility = "private"
	VisibilityPublic          Visibility = "public"
	VisibilityPublicWhitelist Visibility = "public-whitelist"
	VisibilityESP             Visibility = "esp"
	VisibilityESPv2           Visibility = "espv2"
	VisibilityIAP             Visibility = "iap"
	VisibilityApigee          Visibility = "apigee"

	VisibilityUnknown Visibility = ""
)
