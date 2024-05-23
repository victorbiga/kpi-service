package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Create a custom registry
var customRegistry = prometheus.NewRegistry()

// Define gauge metrics with labels
var terraformResourcesCount = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "terraform_resources_count",
		Help: "Current number of Terraform resources.",
	},
	[]string{"organization", "repository"},
)

var cnrmResourcesCount = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "cnrm_resources_count",
		Help: "Current number of CNRM resources.",
	},
	[]string{"organization", "repository"},
)

func init() {
	// Register metrics with the custom registry
	customRegistry.MustRegister(terraformResourcesCount)
	customRegistry.MustRegister(cnrmResourcesCount)
}

func main() {
	// Handle POST requests to update metrics
	http.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		updateMetric(w, r)
	})

	// Expose the custom registry metrics via HTTP
	http.Handle("/metrics", promhttp.HandlerFor(customRegistry, promhttp.HandlerOpts{}))

	// Liveness probe
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// Readiness probe
	http.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ready"))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func updateMetric(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Number       float64 `json:"number"`
		Organization string  `json:"organization"`
		Repository   string  `json:"repository"`
		Metric       string  `json:"metric"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if data.Organization == "" || data.Repository == "" || data.Metric == "" {
		http.Error(w, "Missing 'organization', 'repository', or 'metric' field", http.StatusBadRequest)
		return
	}
	// Log the content of the request body
	log.Printf("Received POST request: %+v\n", data)

	// Choose the appropriate metric to update based on the "metric" field in the request
	var metricToUpdate *prometheus.GaugeVec
	switch data.Metric {
	case "terraform_resources_count":
		metricToUpdate = terraformResourcesCount
	case "cnrm_resources_count":
		metricToUpdate = cnrmResourcesCount
	default:
		http.Error(w, "Invalid 'metric' field value", http.StatusBadRequest)
		return
	}
	metricToUpdate.WithLabelValues(data.Organization, data.Repository).Set(data.Number)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Metric updated"))

	// Log the request and the updated metric
	log.Printf("Updated metric: %s, Organization: %s, Repository: %s, Number: %f\n",
		data.Metric, data.Organization, data.Repository, data.Number)
}
