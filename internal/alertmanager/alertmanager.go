package alertmanager

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Alert struct {
	Labels       map[string]string `json:"labels"`
	Annotations  map[string]string `json:"annotations"`
	StartsAt     time.Time         `json:"startsAt"`
	EndsAt       time.Time         `json:"endsAt"`
	GeneratorURL string            `json:"generatorURL"`
}

func sendAlert(alerts []Alert, alertmanagerURL string) error {
	data, err := json.Marshal(alerts)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", alertmanagerURL, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	log.Printf("Alert sent. Status Code: %d", resp.StatusCode)
	return nil
}

func AlertHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		alertmanagerURL := "http://localhost:9093/api/v2/alerts"

		alerts := []Alert{
			{
				Labels:       map[string]string{"alertname": "TestAlert", "service": "myapp"},
				Annotations:  map[string]string{"summary": "A test alert from myapp."},
				StartsAt:     time.Now(),
				EndsAt:       time.Now().Add(time.Hour),
				GeneratorURL: "http://example.com/graph",
			},
		}
		err := sendAlert(alerts, alertmanagerURL)
		if err != nil {
			http.Error(writer, fmt.Sprintf("Failed to send alert: %v", err), http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte("Alert sent successfully"))
	}
}
