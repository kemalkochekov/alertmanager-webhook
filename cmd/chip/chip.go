package main

import (
	"alartmanagerWebhook/internal/prometheuserver"
	"context"
	"errors"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"os/signal"
	"syscall"
)

func sendMetric() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		randNumber := rand.Intn(10)
		for i := 0; i < randNumber; i++ {
			prometheuserver.OpsProcessed.Inc()
		}
		w.WriteHeader(http.StatusOK)
	}
}
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	router := mux.NewRouter()
	router.HandleFunc("/metric", sendMetric()).Methods(http.MethodGet)
	router.HandleFunc("/welcome", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("welcome chip system"))
	}).Methods(http.MethodGet)

	srv := http.Server{}
	srv.Addr = "localhost:9000"
	srv.Handler = router

	go func() {
		log.Println("The chip server is starting on port 9000")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()

	go func() {
		log.Println("The prometheus server is starting on port 9090")
		if err := prometheuserver.Server(); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()
	log.Println("got interruption signal")
	if err := srv.Shutdown(context.TODO()); err != nil {
		log.Printf("server shutdown returned an err: %v\n", err)
	}

	log.Printf("Graceful Shut down finished")
}
