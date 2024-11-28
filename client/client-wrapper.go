package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ClientErrorHandler func(*http.Response) error

type HttpClient struct {
	BaseURL    string
	Client     *http.Client
	AuthHeader string
}

type RequestOptions struct {
	Method       string
	Path         string
	Payload      interface{}
	Response     interface{}
	Headers      map[string]string
	ErrorHandler ClientErrorHandler
}

func NewHttpClient(baseURL, authHeader string) *HttpClient {
	return &HttpClient{
		BaseURL:    baseURL,
		Client:     &http.Client{},
		AuthHeader: authHeader,
	}
}

func (hc *HttpClient) DoRequest(opts *RequestOptions, errorHandler func(err error)) error {
	fullUrl := hc.BaseURL + opts.Path
	var reqBody io.Reader
	if opts.Payload != nil {
		jsonBody, err := json.Marshal(opts.Payload)
		if err != nil {
			return fmt.Errorf("Error marshaling the request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}
	req, err := http.NewRequest(opts.Method, fullUrl, reqBody)
	if err != nil {
		return fmt.Errorf("Error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", hc.AuthHeader)
	for key, value := range opts.Headers {
		req.Header.Set(key, value)
	}
	resp, err := hc.Client.Do(req)
	if err != nil {
		return fmt.Errorf("Error performing request: %w", err)
	}
	defer resp.Body.Close()
	if opts.ErrorHandler != nil {
    resp.Body.Close()
    errorHandler(err)
    return nil
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Error reading response body: %w", err)
	}
	if opts.Response != nil {
		if err := json.Unmarshal(body, opts.Response); err != nil {
			return fmt.Errorf("Error unmarshaling response: %w", err)
		}
	}
	return nil
}
