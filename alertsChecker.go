package main

import (
	"encoding/json"
	"net/http"
	"time"
)

var (
	alertsURL = "https://prometheus-staging.travix.com/api/v1/alerts"
)

var myClient = &http.Client{Timeout: 10 * time.Second}

func checkAlerts(params BabysitterParams) (bool, error) {
	start := time.Now()
	endgame := start.Add(time.Second * time.Duration(params.WatchTimeSec))

	for {
		alerted, err := wasAlerted(params.PrometheusAlerts)

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

func wasAlerted(alertTypes []string) (bool, error) {

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
