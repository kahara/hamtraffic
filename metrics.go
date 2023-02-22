package hamtraffic

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"net/http"
)

var (
	metrics = make(map[string]*prometheus.CounterVec)
)

func Metrics() {
	metric := promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: PrometheusNamespace,
		Subsystem: "station",
		Name:      "transmissions_total",
		Help:      "Total transmissions",
	}, []string{
		"band",
		"mode",
	})
	metrics["transmissions"] = metric

	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(PrometheusAddrPort, nil); err != nil {
		log.Fatal().Err(err).Str("addrport", PrometheusAddrPort).Msg("Could not expose Prometheus metrics")
	}
}
