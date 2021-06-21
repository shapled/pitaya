package spbridge

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ABCRequest struct {
	BaseRequest
	Name string `json:"name"`
}

type ABCResponse struct {
	BaseResponse
	Name string `json:"name"`
}

func ABCHandler(request Request) (Response, error) {
	req := request.(*ABCRequest)
	if req.Name == "test" {
		return &ABCResponse{Name: "ok"}, nil
	}
	return nil, fmt.Errorf("xxx")
}

func Test_Server(t *testing.T) {
	server := NewServer()
	server.POST("/abc", ABCHandler, &ABCRequest{})
	go server.Start(":10200")
	defer server.Close()
	bodyJson, err := json.Marshal(map[string]interface{}{
		"name": "test",
	})
	assert.Nil(t, err)
	req, err := http.NewRequest(http.MethodPost, "http://localhost:10200/abc", bytes.NewBuffer(bodyJson))
	req.Header.Set("Content-Type", "application/json")
	assert.Nil(t, err)
	resp, err := http.DefaultClient.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	r := &ABCResponse{}
	err = json.Unmarshal(bs, r)
	assert.Nil(t, err)
	assert.Equal(t, "ok", r.Name)
}
