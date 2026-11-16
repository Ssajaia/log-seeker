package main

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
)

type fakeStore struct {
	entries []LogEntry
	got     LogFilter
}

func (f *fakeStore) List(_ context.Context, filter LogFilter) ([]LogEntry, error) {
	f.got = filter
	return f.entries, nil
}

func TestAPIForwardsFiltersAndPagination(t *testing.T) {
	store := &fakeStore{entries: []LogEntry{{ID: 1, Level: "ERROR", Message: "connection timeout"}, {ID: 2, Level: "ERROR", Message: "healthy"}}}
	req := httptest.NewRequest("GET", "/api/logs?level=error&regex=timeout&limit=20&offset=2", nil)
	response := httptest.NewRecorder()
	NewAPI(store).ServeHTTP(response, req)
	if response.Code != 200 {
		t.Fatalf("status %d: %s", response.Code, response.Body)
	}
	var body struct {
		Logs []LogEntry `json:"logs"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if len(body.Logs) != 2 {
		t.Fatalf("unexpected logs: %#v", body.Logs)
	}
	if store.got.Level != "ERROR" || store.got.Limit != 20 || store.got.Offset != 2 {
		t.Fatalf("unexpected filter: %#v", store.got)
	}
}

func TestAPIInvalidRegexAndPagination(t *testing.T) {
	for _, target := range []string{`/api/logs?regex=[`, `/api/logs?limit=0`, `/api/logs?offset=-1`, `/api/logs?regex=` + strings.Repeat("a", 4097)} {
		response := httptest.NewRecorder()
		NewAPI(&fakeStore{}).ServeHTTP(response, httptest.NewRequest("GET", target, nil))
		if response.Code != 400 {
			t.Errorf("%s: status %d", target, response.Code)
		}
		var body map[string]string
		if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil || body["error"] == "" {
			t.Errorf("%s: missing error response", target)
		}
	}
}

func TestAPIMethodNotAllowed(t *testing.T) {
	response := httptest.NewRecorder()
	NewAPI(&fakeStore{}).ServeHTTP(response, httptest.NewRequest("POST", "/api/logs", nil))
	if response.Code != 405 {
		t.Fatalf("status %d, want 405", response.Code)
	}
}
