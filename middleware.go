package stackdriver

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/felixge/httpsnoop"
	"github.com/hellofresh/logging-go/context"
	"github.com/sirupsen/logrus"
)

// LoggingMiddleware is a middleware for writing request logs in a stuctured
// format to stackdriver.
func LoggingMiddleware(log *logrus.Logger) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.New(r.Context()))

			// https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#HttpRequest
			request := &HTTPRequest{
				RequestMethod: r.Method,
				RequestURL:    r.RequestURI,
				RemoteIP:      r.RemoteAddr,
				Referer:       r.Referer(),
				UserAgent:     r.UserAgent(),
				RequestSize:   strconv.FormatInt(r.ContentLength, 10),
			}

			m := httpsnoop.CaptureMetrics(handler, w, r)

			request.Status = strconv.Itoa(m.Code)
			request.Latency = fmt.Sprintf("%.9fs", m.Duration.Seconds())
			request.ResponseSize = strconv.FormatInt(m.Written, 10)

			fields := logrus.Fields{"httpRequest": request}

			// No idea if this works
			traceHeader := r.Header.Get("X-Cloud-Trace-Context")
			if traceHeader != "" {
				fields["trace"] = traceHeader
			}

			log.WithFields(fields).Info("Completed request")
		})
	}
}
