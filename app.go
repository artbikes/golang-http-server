package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	var name, _ = os.Hostname()
	fmt.Fprintf(w, "<h1>This request was processed by host: %s</h1>\n", name)
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func openLogFile(logfile string) {
	if logfile != "" {
		lf, err := os.OpenFile(logfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0640)

		if err != nil {
			log.Fatal("OpenLogfile: os.OpenFile:", err)
		}

		log.SetOutput(lf)
	}
}

func main() {
	currentTime := time.Now().Local().Format("2006-01-02")
	logPath := fmt.Sprintf("/var/log/api/api-%s.log", currentTime)
	httpPort := 8000
	openLogFile(logPath)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	fmt.Fprintf(os.Stdout, "Web Server started. Listening on 0.0.0.0:%v\n", httpPort)
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(fmt.Sprintf(":%d", httpPort), logRequest(http.DefaultServeMux))
	if err != nil {
		log.Fatal(err)
	}
}
