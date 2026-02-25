package httpx

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var _ Client = (*RetryClient)(nil)

type RetryClient struct {
	httpClient    *http.Client
	retryPolicy   RetryPolicy
	backoffPolicy BackoffPolicy
}

func NewRetryClient(opts ...RetryClientOption) *RetryClient {
	httpClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:          100,
			MaxIdleConnsPerHost:   10,
			IdleConnTimeout:       90 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
		},
		Timeout: 15 * time.Second,
	}
	r := &RetryClient{
		httpClient:    httpClient,
		retryPolicy:   StandardRetry(3),
		backoffPolicy: ExponentialBackoffWithJitter(1*time.Second, 10*time.Second),
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

type RetryClientOption func(*RetryClient)

func WithHTTPClient(httpClient *http.Client) RetryClientOption {
	return func(c *RetryClient) {
		c.httpClient = httpClient
	}
}

func WithRetryPolicy(retryPolicy RetryPolicy) RetryClientOption {
	return func(c *RetryClient) {
		c.retryPolicy = retryPolicy
	}
}

func WithBackoffPolicy(backoffPolicy BackoffPolicy) RetryClientOption {
	return func(c *RetryClient) {
		c.backoffPolicy = backoffPolicy
	}
}

func (c *RetryClient) Do(req *http.Request) (*http.Response, error) {
	var bodyBytes []byte
	if req.Body != nil {
		var err error
		bodyBytes, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, fmt.Errorf("reading request body: %w", err)
		}
		req.Body.Close()
	}

	var resp *http.Response
	var err error
	for attempt := 0; ; attempt++ {
		if bodyBytes != nil {
			req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
			req.ContentLength = int64(len(bodyBytes))
			req.GetBody = func() (io.ReadCloser, error) {
				return io.NopCloser(bytes.NewReader(bodyBytes)), nil
			}
		}

		resp, err = c.httpClient.Do(req)

		var respBodyBytes []byte
		if resp != nil {
			respBodyBytes, _ = io.ReadAll(resp.Body)
			resp.Body.Close()
			resp.Body = io.NopCloser(bytes.NewReader(respBodyBytes))
		}

		retry, returnErr := c.retryPolicy(resp, err, attempt)
		if !retry {
			if resp != nil {
				resp.Body = io.NopCloser(bytes.NewReader(respBodyBytes))
			}
			if returnErr != nil {
				return resp, returnErr
			}
			return resp, err
		}

		if returnErr != nil {
			err = returnErr
		}

		time.Sleep(c.backoffPolicy(attempt))
	}
}

func (c *RetryClient) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *RetryClient) Post(url string, contentType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return c.Do(req)
}

func (c *RetryClient) PostForm(url string, data url.Values) (*http.Response, error) {
	return c.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
}
