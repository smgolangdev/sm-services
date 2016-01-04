package main

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"golang.org/x/net/context"
)

type proxymw struct {
	context.Context
	UppercaseEndpoint endpoint.Endpoint
	StringService
}

func proxyingMiddleware(proxyList string, ctx context.Context,
	logger log.Logger) ServiceMiddleware {
	return

}
