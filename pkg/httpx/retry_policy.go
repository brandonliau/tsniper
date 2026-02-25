package httpx

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
)

var (
	redirectsErrorRe     = regexp.MustCompile(`stopped after \d+ redirects\z`)
	schemeErrorRe        = regexp.MustCompile(`unsupported protocol scheme`)
	invalidHeaderErrorRe = regexp.MustCompile(`invalid header`)
	notTrustedErrorRe    = regexp.MustCompile(`certificate is not trusted`)
)

type RetryPolicy func(resp *http.Response, err error, attempt int) (bool, error)

func StandardRetry(maxAttempts int) RetryPolicy {
	return func(resp *http.Response, err error, attempt int) (bool, error) {
		if attempt >= maxAttempts {
			return false, err
		}
		if err != nil {
			if v, ok := err.(*url.Error); ok {
				if redirectsErrorRe.MatchString(v.Error()) {
					return false, v
				}
				if schemeErrorRe.MatchString(v.Error()) {
					return false, v
				}
				if invalidHeaderErrorRe.MatchString(v.Error()) {
					return false, v
				}
				if notTrustedErrorRe.MatchString(v.Error()) {
					return false, v
				}
			}
			return true, nil
		}
		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusRequestTimeout {
			return true, nil
		}
		if resp.StatusCode >= 500 && resp.StatusCode != http.StatusNotImplemented {
			return true, fmt.Errorf("unexpected HTTP status %s", resp.Status)
		}
		return false, nil
	}
}

func NoRetry() RetryPolicy {
	return func(resp *http.Response, err error, attempt int) (bool, error) {
		return false, err
	}
}
