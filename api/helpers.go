package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/sethgrid/pester"
)

// a valid label must be an empty string or consist of alphanumeric characters, '-', '_' or '.', and must start and end with an alphanumeric character (e.g. 'MyValue',  or 'my_value',  or '12345', regex used for validation is '(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?')
func SanitizeLabel(value string) string {

	// Valid label values must be 63 characters or less and must be empty or begin and end with an alphanumeric character ([a-z0-9A-Z])
	// with dashes (-), underscores (_), dots (.), and alphanumerics between.

	// replace @ with -at-
	reg := regexp.MustCompile(`@+`)
	value = reg.ReplaceAllString(value, "-at-")

	// replace all invalid characters with a hyphen
	reg = regexp.MustCompile(`[^a-zA-Z0-9-_.]+`)
	value = reg.ReplaceAllString(value, "-")

	// replace double hyphens with a single one
	value = strings.Replace(value, "--", "-", -1)

	// ensure it starts with an alphanumeric character
	reg = regexp.MustCompile(`^[^a-zA-Z0-9]+`)
	value = reg.ReplaceAllString(value, "")

	// maximize length at 63 characters
	if len(value) > 63 {
		value = value[:63]
	}

	// ensure it ends with an alphanumeric character
	reg = regexp.MustCompile(`[^a-zA-Z0-9]+$`)
	value = reg.ReplaceAllString(value, "")

	return value
}

func SanitizeLabels(labels map[string]string) (sanitizedLabels map[string]string) {
	sanitizedLabels = make(map[string]string, len(labels))
	for k, v := range labels {
		sanitizedLabels[k] = SanitizeLabel(v)
	}
	return
}

func httpRequestHeader(method, url string, headers map[string]string, responseHeader string) string {
	client := pester.New()
	client.MaxRetries = 3
	client.Backoff = pester.ExponentialJitterBackoff
	client.KeepLog = true
	client.Timeout = time.Second * 5
	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		return ""
	}

	for k, v := range headers {
		request.Header.Add(k, v)
	}

	// perform actual request
	response, err := client.Do(request)
	if err != nil {
		return ""
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return ""
	}

	return response.Header.Get(responseHeader)
}

func httpRequestBody(method, url string, headers map[string]string) string {
	client := pester.New()
	client.MaxRetries = 3
	client.Backoff = pester.ExponentialJitterBackoff
	client.KeepLog = true
	client.Timeout = time.Second * 5
	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		return ""
	}

	for k, v := range headers {
		request.Header.Add(k, v)
	}

	// perform actual request
	response, err := client.Do(request)
	if err != nil {
		return ""
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return ""
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return ""
	}

	return string(body)
}

func GetTrimmedDate(date string) (string, error) {
	t, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return "", err
	}
	return fmt.Sprint(t.Format("2006-01-02 15:04:05")), nil
}
