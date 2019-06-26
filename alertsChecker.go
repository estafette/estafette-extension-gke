package main

import (
	"encoding/json"
	"net/http"
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

func checkAlerts(params *Params) (bool, error) {
	start := time.Now()
	alertsURL := getAlertURL(params.Namespace)
	endgame := start.Add(time.Second * time.Duration(params.Babysitter.WatchTimeSec))

	for {
		alerted, err := wasAlerted(params.Babysitter.PrometheusAlerts, alertsURL)

		if alerted || err != nil {
			//TODO: slack message
			return false, err
		}

		if start.After(endgame) {
			return true, nil
		}

		time.Sleep(10 * time.Second)
	}
}

func wasAlerted(alertTypes []string, alertsURL string) (bool, error) {

	alerts := new(alertsResponse)
	err := getJSON(alertsURL, alerts)

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

func getJSON(url string, target interface{}) error {

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", authToken)

	r, err := myClient.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}
