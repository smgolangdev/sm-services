package main

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"golang.org/x/net/context"
	"net/http"
	"os"
	"time"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	httptransport "github.com/go-kit/kit/transport/http"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

func main() {
	ctx := context.Background()
	logger := log.NewLogfmtLogger(os.Stderr)

	//Define metrics to be collected  for this service
	fieldKeys := []string{"method", "error"}
	requestCount := kitprometheus.NewCounter(stdprometheus.CounterOpts{
		Namespace: "my_group",
		Subsystem: "string_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)

	requestLatency := metrics.NewTimeHistogram(time.Microsecond, kitprometheus.NewSummary(stdprometheus.SummaryOpts{
		Namespace: "my_group",
		Subsystem: "string_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in Microseconds.",
	}, fieldKeys))

	countResult := kitprometheus.NewSummary(stdprometheus.SummaryOpts{
		Namespace: "my_group",
		Subsystem: "string_service",
		Name:      "count_result",
		Help:      "The result of eac count method",
	}, []string{})

	//Declare the string service
	var svc StringService
	svc = stringService{}
	svc = loggingMiddleware{logger, svc}
	svc = instrumentingMiddleware{requestCount, requestLatency, countResult, svc}

	//Declare and define a http server that exposes the string service to clients via httptransport
	uppercaseHandler := httptransport.NewServer(
		ctx,
		makeUppercaseEndpoint(svc),
		decodeUppercaseRequest,
		encodeResponse,
	)

	countHandler := httptransport.NewServer(
		ctx,
		makeCountEndpoint(svc),
		decodeCountRequest,
		encodeResponse,
	)

	//Define the content routes
	http.Handle("/uppercase", uppercaseHandler)
	http.Handle("/count", countHandler)
	http.Handle("/metrics", stdprometheus.Handler())
	logger.Log("msg", "HTTP", "addr", ":8080")
	logger.Log("err", http.ListenAndServe(":8080", nil))
}
