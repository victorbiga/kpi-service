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
	http.HandleFunc("/update", logRequest(updateMetric, terraformResourcesCount))

	// Expose the custom registry metrics via HTTP
	http.Handle("/metrics", promhttp.HandlerFor(customRegistry, promhttp.HandlerOpts{}))

	// Liveness probe
	http.HandleFunc("/healthz", logRequest(healthzHandler, nil))

	// Readiness probe
	http.HandleFunc("/readyz", logRequest(readyzHandler, nil))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

// logRequest is a middleware function that logs the request details before passing it to the handler.
func logRequest(next http.HandlerFunc, metric *prometheus.GaugeVec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)
		log.Printf("Request headers: %+v", r.Header)
		if metric != nil {
			// Log additional metric details if available
			log.Printf("Metric name: %s", metric.Desc().String())
		}
		next.ServeHTTP(w, r)
	}
}

// updateMetric updates the specified metric with data from the request body.
func updateMetric(w http.ResponseWriter, r *http.Request, metric *prometheus.GaugeVec) {
	var data struct {
		Number       float64 `json:"number"`
		Organization string  `json:"organization"`
		Repository   string  `json:"repository"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if data.Organization == "" || data.Repository == "" {
		http.Error(w, "Missing 'organization' or 'repository' field", http.StatusBadRequest)
		return
	}
	metric.WithLabelValues(data.Organization, data.Repository).Set(data.Number)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Metric updated"))
}

// healthzHandler handles the liveness probe.
func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

// readyzHandler handles the readiness probe.
func readyzHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the application is ready to serve traffic
	// Implement your readiness checks here
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ready"))
}
