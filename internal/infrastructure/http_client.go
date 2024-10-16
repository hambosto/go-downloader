package infrastructure

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type HTTPClient struct {
	client     *http.Client
	maxRetries int
	retryDelay time.Duration
}

func NewHTTPClient(timeout time.Duration, maxRetries int, retryDelay time.Duration) *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: timeout,
		},
		maxRetries: maxRetries,
		retryDelay: retryDelay,
	}
}

func (c *HTTPClient) GetFileSize(url string) (int64, error) {
	var resp *http.Response
	var err error

	for retry := 0; retry <= c.maxRetries; retry++ {
		resp, err = c.client.Get(url)
		if err == nil {
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return resp.ContentLength, nil
			}
			err = fmt.Errorf("unexpected status: %s", resp.Status)
		}

		if retry < c.maxRetries {
			time.Sleep(c.retryDelay)
		}
	}

	return 0, fmt.Errorf("failed to get file size after %d retries: %w", c.maxRetries, err)
}

func (c *HTTPClient) DownloadChunk(url string, offset, size int64) (io.ReadCloser, error) {
	var resp *http.Response
	var err error

	for retry := 0; retry <= c.maxRetries; retry++ {
		var req *http.Request
		req, err = http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return nil, fmt.Errorf("creating request: %w", err)
		}

		rangeHeader := fmt.Sprintf("bytes=%d-%d", offset, offset+size-1)
		req.Header.Set("Range", rangeHeader)

		resp, err = c.client.Do(req)
		if err == nil && resp.StatusCode == http.StatusPartialContent {
			return resp.Body, nil
		}

		if resp != nil {
			resp.Body.Close()
		}

		if retry < c.maxRetries {
			time.Sleep(c.retryDelay)
		}
	}

	return nil, fmt.Errorf("failed to download chunk after %d retries: %w", c.maxRetries, err)
}
