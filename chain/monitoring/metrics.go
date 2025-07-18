package monitoring

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var collateralGauge = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "usdWz_collateral",
	Help: "Current M0 collateral backing",
})

// StartServer exposes Prometheus metrics on addr.
func StartServer(addr string) error {
	prometheus.MustRegister(collateralGauge)
	http.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(addr, nil)
}

// SetCollateral sets the gauge value.
func SetCollateral(v float64) { collateralGauge.Set(v) }
