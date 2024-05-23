package main

import (
    "encoding/json"
    "net/http"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "log"
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
        updateMetric(w, r, terraformResourcesCount)
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
        // Check if the application is ready to serve traffic
        // Implement your readiness checks here
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("ready"))
    })

    log.Fatal(http.ListenAndServe(":8080", nil))
}

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
