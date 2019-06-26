package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

var myClient = &http.Client{Timeout: 10 * time.Second}

func getAlertURL(namespace string) string {
	switch namespace {
	case "production":
		return "https://prometheus-production.travix.com/api/v1/alerts"
	default:
		return "https://prometheus-staging.travix.com/api/v1/alerts"
	}
}

func checkAlerts(params Params) (bool, error) {
	alertsURL := getAlertURL(params.Namespace)
	endgame := time.Now().Add(time.Second * time.Duration(params.Babysitter.WatchTimeSec))

	for {
		alerted, err := wasAlerted(params.Babysitter.PrometheusAlerts, alertsURL, params.Babysitter.PrometheusToken)

		if alerted || err != nil {
			alertedStr := strconv.FormatBool(alerted)
			logInfo("Checking alerts failed. Alerted: "+alertedStr+" With errors: ", err)

			return false, err
		}

		if time.Now().After(endgame) {
			logInfo("Checking alerts passed, no errors found")
			return true, nil
		}

		time.Sleep(10 * time.Second)
	}
}

func wasAlerted(alertTypes []string, alertsURL string, authToken string) (bool, error) {

	alerts := new(alertsResponse)
	err := getJSON(alertsURL, alerts, authToken)

	if err != nil {
		return false, err
	}

	hash := make(map[string]bool)

	for _, alertType := range alertTypes {
		hash[alertType] = true
	}

	for idx := range alerts.Data.Alerts {
		if hash[alerts.Data.Alerts[idx].Labels.Alertname] {
			return true, nil
		}
	}

	return false, nil
}

func getJSON(url string, target interface{}, authToken string) error {

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", authToken)

	r, err := myClient.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}
