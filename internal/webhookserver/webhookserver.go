package webhookserver

import (
	"io"
	"log"
	"net/http"
)

func WebhookHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		log.Printf("Called webhook handler")
		body, err := io.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, "Error reading body", http.StatusInternalServerError)
			return
		}
		defer request.Body.Close()
		log.Printf("Received alert: %s\n", string(body))
	}
}
