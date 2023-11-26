package api

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
)

var ErrInvalidStatusCode = errors.New("invalid status code")

// GetRedirect returns the final URL after first redirection.
func GetRedirect(url string) (string, error) {
	const op = "api.GetRedirect"

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // stop after 1st redirect
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusFound {
		return "", fmt.Errorf("%s: %w: %d", op, ErrInvalidStatusCode, resp.StatusCode)
	}

	return resp.Header.Get("Location"), nil
}

func GetDelete(url string) (*http.Response, error) {
	const op = "api.GetDelete"
	r, err := http.NewRequest(http.MethodDelete, url, bytes.NewReader([]byte{}))
	if err != nil {
		return nil, fmt.Errorf("new delete request for url %s: %w", url, err)
	}
	defer r.Body.Close()
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, fmt.Errorf("do request for url %s: %w", url, err)
	}
	return resp, nil
}
