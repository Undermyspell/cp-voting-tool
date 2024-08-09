package httputils

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func Get[T any](url string, headers map[string]string) (*T, int) {
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		logrus.Error(err)
		return nil, http.StatusInternalServerError
	}

	for k, header := range headers {
		req.Header.Set(k, header)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	// Create an HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusInternalServerError
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, http.StatusInternalServerError
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusInternalServerError
	}

	var res T
	err = json.Unmarshal(body, &res)
	if err != nil {
		logrus.Error(err)
		return nil, http.StatusInternalServerError
	}

	return &res, resp.StatusCode
}
