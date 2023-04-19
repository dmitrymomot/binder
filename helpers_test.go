package binder_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

// new json request with body and headers
func newJSONRequest(method, url string, body interface{}, headers map[string]string) (*http.Request, error) {
	var buf *bytes.Buffer
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		buf = bytes.NewBuffer(b)
	}

	var req *http.Request
	var err error
	if buf != nil {
		req, err = http.NewRequest(method, url, buf)
	} else {
		req, err = http.NewRequest(method, url, nil)
	}
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return req, nil
}

// new query request with headers
func newQueryRequest(method, url string, query map[string]interface{}, headers map[string]string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	q := req.URL.Query()
	for key, value := range query {
		q.Add(key, fmt.Sprintf("%v", value))
	}
	req.URL.RawQuery = q.Encode()

	return req, nil
}

// new form-urlencoded request with body and headers
func newFormRequest(method, endpoint string, body map[string]interface{}, headers map[string]string) (*http.Request, error) {
	var reqBody io.Reader
	if body != nil {
		q := url.Values{}
		for key, value := range body {
			q.Add(key, fmt.Sprintf("%v", value))
		}
		reqBody = strings.NewReader(q.Encode())
	}

	req, err := http.NewRequest(method, endpoint, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return req, nil
}

// new multipart request with body and headers
func newMultipartRequest(method, url string, body map[string]interface{}, headers map[string]string) (*http.Request, error) {
	var contentType string = "multipart/form-data"
	var buf bytes.Buffer
	if len(body) > 0 {
		w := multipart.NewWriter(&buf)
		for key, value := range body {
			if _, ok := value.(io.Reader); ok {
				if fw, err := w.CreateFormFile(key, key); err != nil {
					return nil, err
				} else if _, err := io.Copy(fw, value.(io.Reader)); err != nil {
					return nil, err
				}
			} else {
				if err := w.WriteField(key, fmt.Sprintf("%v", value)); err != nil {
					return nil, err
				}
			}
		}
		if err := w.Close(); err != nil {
			return nil, err
		}
		contentType = w.FormDataContentType()
	}

	req, err := http.NewRequest(method, url, &buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return req, nil
}
