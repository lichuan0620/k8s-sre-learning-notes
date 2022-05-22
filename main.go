package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

func main() {
	const listenAddress = "0.0.0.0:8080"
	http.Handle("/metrics", promhttp.Handler())
	logrus.Infof("listen and serve on %s", listenAddress)
	if err := http.ListenAndServe(listenAddress, http.DefaultServeMux); err != http.ErrServerClosed {
		logrus.Fatalln(err)
	}
}
