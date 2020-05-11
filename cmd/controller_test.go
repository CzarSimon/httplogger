package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CzarSimon/httplogger/internal/models"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestLog(t *testing.T) {
	assert := assert.New(t)
	s := server()

	errorEvent := models.Event{
		AppName:    "TEST_APP",
		Version:    "1.2.0",
		SessionID:  "some-session-id",
		ClientID:   "some-client-id",
		Message:    "A wild event occured.",
		Stacktrace: "this is a stacktrace.",
		Level:      "error",
	}

	req := createTestRequest("/v1/logs", http.MethodPost, errorEvent)
	res := performTestRequest(s.Handler, req)
	assert.Equal(http.StatusOK, res.Code)

	infoEvent := models.Event{
		AppName:   "TEST_APP",
		Version:   "1.2.0",
		SessionID: "some-session-id",
		ClientID:  "some-client-id",
		Message:   "A wild event occured.",
		Level:     "info",
	}

	req = createTestRequest("/v1/logs", http.MethodPost, infoEvent)
	res = performTestRequest(s.Handler, req)
	assert.Equal(http.StatusOK, res.Code)

	unsupportedLogLevel := models.Event{
		AppName:   "TEST_APP",
		Version:   "1.2.0",
		SessionID: "some-session-id",
		ClientID:  "some-client-id",
		Message:   "A wild event occured.",
		Level:     "panic",
	}

	req = createTestRequest("/v1/logs", http.MethodPost, unsupportedLogLevel)
	res = performTestRequest(s.Handler, req)
	assert.Equal(http.StatusBadRequest, res.Code)

	unsupportedEvent := []string{
		"bla",
		"blabla",
	}

	req = createTestRequest("/v1/logs", http.MethodPost, unsupportedEvent)
	res = performTestRequest(s.Handler, req)
	assert.Equal(http.StatusBadRequest, res.Code)
}

func TestLog_BadContentType(t *testing.T) {
	assert := assert.New(t)
	s := server()

	errorEvent := models.Event{
		AppName:    "TEST_APP",
		Version:    "1.2.0",
		SessionID:  "some-session-id",
		ClientID:   "some-client-id",
		Message:    "A wild event occured.",
		Stacktrace: "this is a stacktrace.",
		Level:      "error",
	}

	req := createTestRequest("/v1/logs", http.MethodPost, errorEvent)
	req.Header.Del("Content-Type")
	res := performTestRequest(s.Handler, req)
	assert.Equal(http.StatusUnsupportedMediaType, res.Code)

	req = createTestRequest("/v1/logs", http.MethodPost, errorEvent)
	req.Header.Set("Content-Type", "application/xml")
	res = performTestRequest(s.Handler, req)
	assert.Equal(http.StatusUnsupportedMediaType, res.Code)

	req = createTestRequest("/v1/logs", http.MethodPost, errorEvent)
	req.Header.Set("Content-Type", "text/plain")
	res = performTestRequest(s.Handler, req)
	assert.Equal(http.StatusUnsupportedMediaType, res.Code)

	req = createTestRequest("/v1/logs", http.MethodPost, errorEvent)
	req.Header.Set("Content-Type", "text/html")
	res = performTestRequest(s.Handler, req)
	assert.Equal(http.StatusUnsupportedMediaType, res.Code)
}

func TestCheckHealthAndMetrics(t *testing.T) {
	assert := assert.New(t)

	s := server()
	req := createTestRequest("/health", http.MethodGet, nil)
	res := performTestRequest(s.Handler, req)
	assert.Equal(http.StatusOK, res.Code)

	req = createTestRequest("/metrics", http.MethodGet, nil)
	res = performTestRequest(s.Handler, req)
	assert.Equal(http.StatusOK, res.Code)
}

func performTestRequest(r http.Handler, req *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func createTestRequest(route, method string, body interface{}) *http.Request {
	var reqBody io.Reader
	if body != nil {
		bytesBody, err := json.Marshal(body)
		if err != nil {
			log.Fatal("Failed to marshal body", zap.Error(err))
		}
		reqBody = bytes.NewBuffer(bytesBody)
	}

	req, err := http.NewRequest(method, route, reqBody)
	if err != nil {
		log.Fatal("Failed to create request", zap.Error(err))
	}

	req.Header.Set("Content-Type", "application/json")
	return req
}
