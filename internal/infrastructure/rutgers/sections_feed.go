package rutgers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"tsniper/internal/application/ports"
	"tsniper/internal/domain/scope"

	"tsniper/pkg/httpx"
)

var _ ports.SectionsFeed = (*sectionsFeed)(nil)

const sectionsURL = "https://classes.rutgers.edu//soc/api/openSections.json"

type sectionsFeed struct {
	client *httpx.Client
}

func NewSectionsFeed() *sectionsFeed {
	client := httpx.NewClient(
		httpx.WithRetryPolicy(httpx.NoRetry()),
		httpx.WithBackoffPolicy(httpx.NoBackoff()),
	)
	return &sectionsFeed{
		client: client,
	}
}

func (f *sectionsFeed) FetchOpenSections(scp scope.AcademicScope) ([]string, error) {
	params := url.Values{}
	params.Add("campus", string(scp.Campus))
	params.Add("term", string(scp.Term))
	params.Add("year", scp.Year)

	fullURL := fmt.Sprintf("%s?%s", sectionsURL, params.Encode())
	resp, err := f.client.Get(fullURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("sections feed returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var openSections []string
	if err := json.Unmarshal(body, &openSections); err != nil {
		return nil, err
	}

	return openSections, nil
}
