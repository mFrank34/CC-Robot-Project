package router

import (
	"bytes"
	"io"
	"net/http"
	"time"
)

func readBody(resp *http.Response) ([]byte, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func Get(host, subroute string) ([]byte, error) {
	resp, err := http.Get(host + subroute)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return readBody(resp)
}

func Post(host, subroute string, payload []byte) ([]byte, int, error) {
	url := host + subroute
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	body, err := readBody(resp)
	return body, resp.StatusCode, err
}
