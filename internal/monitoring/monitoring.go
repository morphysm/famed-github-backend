package monitoring

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Monitoring struct {
	Handler http.Handler
}

// just an example
func recordMetrics() {
	go func() {
		for {
			opsProcessed.Inc()
			time.Sleep(2 * time.Second)
		}
	}()
}

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "famed_processed",
		Help: "The total number of processed events",
	})
)

func NewMonitoring() *Monitoring {
	recordMetrics()

	return &Monitoring{Handler: promhttp.Handler()}
}
