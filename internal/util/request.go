package util

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type Request struct {
	URL      string
	Data     any
	ProxyURL string
	Headers  map[string]string
}

func (req *Request) sendRequest(method string) ([]byte, error) {
	data, err := json.Marshal(req.Data)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	if req.ProxyURL != "" {
		proxyURL, err := url.Parse(req.ProxyURL)
		if err != nil {
			return nil, err
		}
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
		client.Transport = transport
	}
	reqBody := bytes.NewBuffer(data)
	httpReq, err := http.NewRequest(method, req.URL, reqBody)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return respBody, nil
}

func (req *Request) POST() ([]byte, error) {
	return req.sendRequest("POST")
}

func (req *Request) GET() ([]byte, error) {
	return req.sendRequest("GET")
}
