package main

import (
	"alartmanagerWebhook/internal/webhookserver"
	"context"
	"errors"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	router := mux.NewRouter()
	router.HandleFunc("/", webhookserver.WebhookHandler()).Methods(http.MethodPost)
	router.HandleFunc("/webhook", webhookserver.WebhookHandler()).Methods(http.MethodPost)
	router.HandleFunc("", webhookserver.WebhookHandler()).Methods(http.MethodPost)
	router.HandleFunc("/welcome", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("welcome webhook alert system"))
	}).Methods(http.MethodGet)

	srv := http.Server{}
	srv.Addr = ":5001"
	srv.Handler = router

	go func() {
		log.Println("The server is starting on port 5001")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen and serve returned err: %v", err)
		}

	}()

	<-ctx.Done()
	log.Println("got interruption signal")
	if err := srv.Shutdown(context.TODO()); err != nil {
		log.Printf("server shutdown returned an err: %v\n", err)
	}

	log.Printf("Graceful Shut down finished")
}
