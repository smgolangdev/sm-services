package main

import (
	"fmt"
	"time"

	"github.com/go-kit/kit/metrics"
)

type instrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.TimeHistogram
	countResult    metrics.Histogram
	StringService
}

func (mw instrumentingMiddleware) Uppercase(s string) (output string, err error) {
	defer func(begin time.Time) {
		methodField := metrics.Field{Key: "method", Value: "Uppercase"}
		errField := metrics.Field{Key: "error", Value: fmt.Sprintf("%v", err)}
		mw.requestCount.With(methodField).With(errField).Add(1)
		mw.requestLatency.With(methodField).With(errField).Observe(time.Since(begin))
	}(time.Now())

	output, err = mw.StringService.Uppercase(s)
	return
}

func (mw instrumentingMiddleware) Count(s string) (c int) {
	defer func(begin time.Time) {
		methodField := metrics.Field{Key: "method", Value: "Count"}
		errField := metrics.Field{Key: "error", Value: fmt.Sprintf("%v", error(nil))}
		mw.requestCount.With(methodField).With(errField).Add(1)
		mw.requestLatency.With(methodField).With(errField).Observe(time.Since(begin))
	}(time.Now())

	c = mw.StringService.Count(s)
	return
}
