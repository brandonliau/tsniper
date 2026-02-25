package httpx

import (
	"io"
	"net/http"
)

func ReadBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func ConsumeBody(resp *http.Response) {
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)
}
