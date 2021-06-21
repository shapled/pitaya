package main

import (
	"errors"
	"log"

	"github.com/shapled/spbridge"
)

type ABCRequest struct {
	spbridge.BaseRequest
	Name string `json:"name"`
}

type ABCResponse struct {
	spbridge.BaseResponse
	Name string `json:"name"`
}

func ABCHandler(request spbridge.Request) (spbridge.Response, error) {
	req := request.(*ABCRequest)
	if req.Name == "test" {
		return &ABCResponse{Name: "ok"}, nil
	}
	return nil, errors.New("xxx")
}

func main() {
	server := spbridge.NewServer()
	server.POST("/abc", ABCHandler, &ABCRequest{})
	log.Fatal(server.Start(":10200"))
}
