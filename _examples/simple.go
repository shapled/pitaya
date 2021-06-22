package main

import (
	"errors"
	"log"

	"github.com/shapled/pitaya"
)

type ABCRequest struct {
	pitaya.BaseRequest
	Name string `json:"name"`
}

type ABCResponse struct {
	pitaya.BaseResponse
	Name string `json:"name"`
}

func ABCHandler(request pitaya.Request) (pitaya.Response, error) {
	req := request.(*ABCRequest)
	if req.Name == "test" {
		return &ABCResponse{Name: "ok"}, nil
	}
	return nil, errors.New("xxx")
}

func main() {
	server := pitaya.NewServer()
	server.POST("/abc", ABCHandler, &ABCRequest{})
	log.Fatal(server.Start(":10200"))
}
