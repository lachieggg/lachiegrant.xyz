package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHealthPublicRevealsNothing(t *testing.T) {
	rec := httptest.NewRecorder()
	healthHandler(rec, httptest.NewRequest(http.MethodGet, "/health", nil))
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"status":"ok"}`, rec.Body.String())
}

func TestHealthDetailsFailsClosedWithoutCredentials(t *testing.T) {
	t.Setenv("HEALTH_USER", "")
	t.Setenv("HEALTH_PASS", "")
	rec := httptest.NewRecorder()
	healthDetailsHandler(rec, httptest.NewRequest(http.MethodGet, "/health/details", nil))
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestHealthDetailsRequiresAuth(t *testing.T) {
	t.Setenv("HEALTH_USER", "ops")
	t.Setenv("HEALTH_PASS", "s3cret")

	rec := httptest.NewRecorder()
	healthDetailsHandler(rec, httptest.NewRequest(http.MethodGet, "/health/details", nil))
	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	req := httptest.NewRequest(http.MethodGet, "/health/details", nil)
	req.SetBasicAuth("ops", "wrong")
	rec = httptest.NewRecorder()
	healthDetailsHandler(rec, req)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	req = httptest.NewRequest(http.MethodGet, "/health/details", nil)
	req.SetBasicAuth("ops", "s3cret")
	rec = httptest.NewRecorder()
	healthDetailsHandler(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "uptime_seconds")
}

func TestRequestSummaryExcludesHealthAndAverages(t *testing.T) {
	reqStats.Lock()
	reqStats.recent, reqStats.next, reqStats.total, reqStats.errors = nil, 0, 0, 0
	reqStats.Unlock()

	recordRequest("/", 200, 10*time.Millisecond)
	recordRequest("/blog", 500, 30*time.Millisecond)
	recordRequest("/health/details", 200, time.Second)

	got := requestSummary()
	assert.Equal(t, 2, got["total"])
	assert.Equal(t, 1, got["errors_5xx"])
	assert.InDelta(t, 20.0, got["avg_ms"], 0.01)
}
