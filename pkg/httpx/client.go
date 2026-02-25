package httpx

import (
	"io"
	"net/http"
	"net/url"
)

type Client interface {
	Do(req *http.Request) (*http.Response, error)
	Get(url string) (*http.Response, error)
	Post(url string, contentType string, body io.Reader) (*http.Response, error)
	PostForm(url string, data url.Values) (*http.Response, error)
}
