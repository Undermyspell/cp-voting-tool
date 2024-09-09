package httputils

import (
	"bytes"
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
		return nil, resp.StatusCode
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

func Post(url string, headers map[string]string, body any) int {

	payload, err := json.Marshal(body)

	if err != nil {
		logrus.Error(err)
		return http.StatusInternalServerError
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))

	if err != nil {
		logrus.Error(err)
		return http.StatusInternalServerError
	}

	for k, header := range headers {
		req.Header.Set(k, header)
	}
	req.Header.Set("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	// Create an HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Error(err)
		return http.StatusInternalServerError
	}
	defer resp.Body.Close()

	return resp.StatusCode
}

func Put(url string, headers map[string]string, body any) int {

	payload, err := json.Marshal(body)

	if err != nil {
		logrus.Error(err)
		return http.StatusInternalServerError
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(payload))

	if err != nil {
		logrus.Error(err)
		return http.StatusInternalServerError
	}

	for k, header := range headers {
		req.Header.Set(k, header)
	}
	req.Header.Set("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	// Create an HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Error(err)
		return http.StatusInternalServerError
	}
	defer resp.Body.Close()

	return resp.StatusCode
}

func Delete(url string, headers map[string]string) int {
	req, err := http.NewRequest("DELETE", url, nil)

	if err != nil {
		logrus.Error(err)
		return http.StatusInternalServerError
	}

	for k, header := range headers {
		req.Header.Set(k, header)
	}
	req.Header.Set("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	// Create an HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Error(err)
		return http.StatusInternalServerError
	}
	defer resp.Body.Close()

	return resp.StatusCode
}
