package negronistatsd

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/peterbourgon/g2s"
)

// Middleware is the middleware handler that sends the status code and response time to StatsD
type Middleware struct {
	Client g2s.Statter
	Prefix string
}

// NewMiddleware returns *Middleware
func NewMiddleware(server string, prefix string) *Middleware {
	statsdclient, err := g2s.Dial("udp", server)
	if err != nil {
		log.Println(err)
	}
	return &Middleware{statsdclient, prefix}
}

func (l *Middleware) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()

	next(rw, r)

	res := rw.(negroni.ResponseWriter)
	responseTime := time.Since(start)

	// Statsd
	timeMetricPath := strings.Join([]string{l.Prefix, "response_time"}, ".")
	go l.Client.Timing(1.0, timeMetricPath, responseTime)

	statusMetricPath := strings.Join([]string{l.Prefix, strconv.Itoa(res.Status())}, ".")
	go l.Client.Counter(1.0, statusMetricPath, 1)

}
