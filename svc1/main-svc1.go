package main

import (
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"strings"

	httptransport "github.com/go-kit/kit/transport/http"
)

//Stringservice provides operations on stringi service
type StringService interface {
	Uppercase(string) (string, error)
	Count(string) int
}

//Implementing the Stringservice interface
type stringService struct{}

var ErrEmpty = errors.New("Empty string")

func (stringService) Uppercase(s string) (string, error) {
	if s == "" {
		return "", ErrEmpty
	}
	return strings.ToUpper(s), nil
}

func (stringService) Count(s string) int {
	return len(s)
}

//Implement requests and responses
//In Go kit, the primary messaging pattern is RPC. So, every method
//in our interface will be modeled as a remote procedure call. For each method,
//we define request and response structs, capturing all of the input and output
//parameters respectively.
type uppercaseRequest struct {
	S string `json:"s"`
}

type uppercaseResponse struct {
	V   string `json:"v"`
	Err string `json:"error omitempty"`
}

type countRequest struct {
	S string `json:"s"`
}

type countResponse struct {
	V int `json"v"`
}

//Now define go endpoints via go-kit abstraction called an Endpoint
func makeUppercaseEndpoint(svc stringService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(uppercaseRequest)
		v, err := svc.Uppercase(req.S)
		if err != nil {
			return uppercaseResponse{v, err.Error()}, nil
		}
		return uppercaseResponse{v, ""}, nil
	}

}

func makeCountEndpoint(svc stringService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(countRequest)
		v := svc.Count(req.S)
		return countResponse{v}, nil
	}
}

//Main
func main() {
	ctx := context.Background()
	svc := stringService{}

	uppercaseHandler := httptransport.NewServer(ctx,
		makeUppercaseEndpoint(svc),
		decodeUppercaseRequest,
		encodeResponse)

	countHandler := httptransport.NewServer(
		ctx, makeCountEndpoint(svc),
		decodeCountRequest,
		encodeResponse)

	http.Handle("/uppercase", uppercaseHandler)
	http.Handle("/count", countHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func decodeUppercaseRequest(r *http.Request) (interface{}, error) {
	var req uppercaseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func decodeCountRequest(r *http.Request) (interface{}, error) {
	var req countRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil

}

func encodeResponse(w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
