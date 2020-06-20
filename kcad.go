package main

import (
	"flag"
	"fmt"
	"github.com/NeoEvolutionAI/k-cad/packages/k8s"
	"log"
	"net/http"
)

func main() {
	var PORT string
	flag.StringVar(&PORT, "port", "3000", "define a port where you would like k-cad to listen on")
	flag.Parse()

	k8s.Initialize()

	log.Printf("K-cad is listening on port: %s\n", PORT)
	http.HandleFunc("/metrics", metricsHandler)

	err := http.ListenAndServe(":"+PORT, nil)
	if err != nil {
		log.Fatal("Couldn't start server", err)
	}
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	metricList, errorsList := k8s.QueryCadvisor()

	if len(errorsList) != 0 {
		fmt.Fprint(w, errorsList)
		return
	}

	for _, m := range metricList {
		_, err := fmt.Fprintf(w, m)
		if err != nil {
			fmt.Fprint(w, "Error writing response", err)
			return
		}
	}

	fmt.Fprintf(w, "")
}
