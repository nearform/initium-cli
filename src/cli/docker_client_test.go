package cli

import (
	"bytes"
	"github.com/docker/docker/client"
	"io"
	"net/http"
)

type dockerApiRequest struct {
	httpMethod string
	url        string
}

type transportFunc func(*http.Request) (*http.Response, error)

func (tf transportFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return tf(req)
}

func newAlwaysOkMockDockerClient(apiRequests *[]dockerApiRequest) (client.Client, error) {
	handler := func(request *http.Request) (*http.Response, error) {
		*apiRequests = append(*apiRequests, dockerApiRequest{
			httpMethod: request.Method,
			url:        request.URL.Path,
		})

		header := http.Header{}
		header.Set("Content-Type", "application/json")

		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte("{}"))),
			Header:     header,
		}, nil
	}
	mockClient, err := client.NewClientWithOpts(client.WithHTTPClient(&http.Client{
		Transport: transportFunc(handler),
	}))
	if err != nil {
		return client.Client{}, err
	}

	return *mockClient, nil
}
