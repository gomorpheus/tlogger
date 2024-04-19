package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func writeLog(log string) error {
	f, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	defer f.Close()
	if _, err = f.WriteString(log + "\n"); err != nil {
		return err
	}
	return nil
}

func middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// middleware operations on the request here
		// ignore favicon
		if r.URL.String() != "/favicon.ico" {
			l := fmt.Sprintf("referrer=%s remote-ip=%s user-agent=%s url=%s host=%s method=%s", r.Referer(), r.RemoteAddr, r.UserAgent(), r.URL.String(), r.Host, r.Method)
			if err := writeLog(l); err != nil {
				log.Println("error writing log entry:", err)
			}
			log.Println(l)
			log.Println("---------")
			log.Printf("%+v\n\n", r)
		}
		h.ServeHTTP(w, r)
	})
}

func main() {
	mux := mux.NewRouter()

	// default "everything" handler
	mux.PathPrefix("/").HandlerFunc(handler).Methods("GET")

	srv := &http.Server{
		Addr:         ":9090",
		Handler:      middleware(mux),
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	// start server
	log.Println("Starting tlogger server listening on port 9090")
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}
