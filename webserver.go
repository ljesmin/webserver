package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const waitForPath = "/waitfor/"

func handler(w http.ResponseWriter, r *http.Request) {
	val, ok := os.LookupEnv("RESPONSE")
	if !ok {
		fmt.Fprintf(w, "Hello, World %q %q\n", r.Host, r.URL.Path)
	} else {
		fmt.Fprintf(w, "%s %q %q\n", val, r.Host, r.URL.Path)
	}
	val, ok = os.LookupEnv("DELAYMS")
	if ok {
		delaytime, err := strconv.ParseUint(val, 10, 32)
		if err == nil {
			time.Sleep(time.Duration(delaytime) * time.Millisecond)
		}
	}
	fmt.Fprintf(w, "*** Headers ***\n")
	for name, headers := range r.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
	body, err := ioutil.ReadAll(r.Body)
	if err == nil {
		fmt.Println(string(body))
	}
}

func healthcheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "{\"healthy\":true}\n")
}

func waitFor(w http.ResponseWriter, r *http.Request) {
	var slug string
	if strings.HasPrefix(r.URL.Path, waitForPath) {
		slug = r.URL.Path[len(waitForPath):]
		delaytime, err := strconv.ParseUint(slug, 10, 32)
		if err == nil {
			time.Sleep(time.Duration(delaytime) * time.Millisecond)
		}
	}
	w.Write([]byte("Waited\n"))
}

func prometheusMetrics(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "http_requests_meie{method=\"post\",code=\"200\"} 1027 1395066363000\nmetric_without_timestamp_and_labels 12.47\n")
}

func main() {
	val, ok := os.LookupEnv("PORT")

	NormalServer := http.NewServeMux()
	NormalServer.HandleFunc("/", handler)

	NormalServer.HandleFunc("/status", healthcheck)
	NormalServer.HandleFunc(waitForPath, waitFor)
	NormalServer.Handle("/metrics", promhttp.Handler())

	if !ok {
		http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, NormalServer))
	} else {
		http.ListenAndServe(":"+val, handlers.LoggingHandler(os.Stdout, NormalServer))
	}
}
