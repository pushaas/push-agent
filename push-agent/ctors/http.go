package ctors

import (
	"net/http"
	"time"

	"github.com/imroc/req"
)

func NewHttpClient() *http.Client {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	return &client
}

func NewReqHttpClient(baseClient *http.Client) *req.Req {
	client := req.New()
	client.SetClient(baseClient)
	return client
}
