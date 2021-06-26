package report_portal

import (
	"context"
	"crypto/tls"
	"io"
	"net/http"

	"github.com/flacatus/che-inspector/pkg/api"
)

type API struct {
	httpClient   *http.Client
	reportPortal *api.ReportPortal
}

func NewReportPortalClient(reportPortal *api.ReportPortal) *API {
	api := API{
		reportPortal: reportPortal,
	}
	api.httpClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	return &api
}

func (c *API) Do(req *http.Request) (*http.Response, error) {
	res, err := c.httpClient.Do(req)
	return res, err
}

func (c *API) Post(ctx context.Context, path string, contentType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", c.reportPortal.BaseUrl+path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", "bearer "+c.reportPortal.Token)
	return c.Do(req)
}
