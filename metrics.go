package hamtraffic

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"net/http"
)

var metrics = make(map[string]*prometheus.CounterVec)

func Metrics() {
	metrics["transmissions"] = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: PrometheusNamespace,
		Subsystem: "station",
		Name:      "transmissions_total",
		Help:      "Total transmissions",
	}, []string{
		"band",
		"mode",
		"callsign",
	})

	metrics["receptions"] = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: PrometheusNamespace,
		Subsystem: "station",
		Name:      "receptions_total",
		Help:      "Total receptions",
	}, []string{
		"band",
		"mode",
		"callsign",
	})

	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(PrometheusAddrPort, nil); err != nil {
		log.Fatal().Err(err).Str("addrport", PrometheusAddrPort).Msg("Could not expose Prometheus metrics")
	}
}
