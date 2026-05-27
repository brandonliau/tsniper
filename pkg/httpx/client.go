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

type Client struct {
	httpClient    *http.Client
	retryPolicy   RetryPolicy
	backoffPolicy BackoffPolicy
}

func NewClient(opts ...ClientOption) *Client {
	httpClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:          100,
			MaxIdleConnsPerHost:   10,
			IdleConnTimeout:       90 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
		},
		Timeout: 15 * time.Second,
	}
	r := &Client{
		httpClient:    httpClient,
		retryPolicy:   StandardRetry(3),
		backoffPolicy: ExponentialBackoffWithJitter(1*time.Second, 10*time.Second),
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
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

func (c *Client) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *Client) Post(url string, contentType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return c.Do(req)
}

func (c *Client) PostForm(url string, data url.Values) (*http.Response, error) {
	return c.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
}

type ClientOption func(*Client)

func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

func WithRetryPolicy(retryPolicy RetryPolicy) ClientOption {
	return func(c *Client) {
		c.retryPolicy = retryPolicy
	}
}

func WithBackoffPolicy(backoffPolicy BackoffPolicy) ClientOption {
	return func(c *Client) {
		c.backoffPolicy = backoffPolicy
	}
}
