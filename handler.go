package main

import (
	"context"
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type LogLister interface {
	List(ctx context.Context, filter LogFilter) ([]LogEntry, error)
}

type API struct{ store LogLister }

func NewAPI(store LogLister) *API { return &API{store: store} }

func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	q := r.URL.Query()
	level := strings.ToUpper(strings.TrimSpace(q.Get("level")))
	if level != "" && len(level) > 32 {
		writeError(w, http.StatusBadRequest, "level must be at most 32 characters")
		return
	}
	regexText := q.Get("regex")
	if regexText != "" {
		if len(regexText) > 4096 {
			writeError(w, http.StatusBadRequest, "regex must be at most 4096 characters")
			return
		}
		if _, err := regexp.Compile(regexText); err != nil {
			writeError(w, http.StatusBadRequest, "invalid regex: "+err.Error())
			return
		}
	}
	limit, err := queryInt(q.Get("limit"), 100)
	if err != nil || limit < 1 || limit > 500 {
		writeError(w, http.StatusBadRequest, "limit must be an integer from 1 to 500")
		return
	}
	offset, err := queryInt(q.Get("offset"), 0)
	if err != nil || offset < 0 {
		writeError(w, http.StatusBadRequest, "offset must be a non-negative integer")
		return
	}
	logs, err := a.store.List(r.Context(), LogFilter{Level: level, Regex: regexText, Limit: limit, Offset: offset})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to retrieve logs")
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(map[string]any{"logs": logs, "limit": limit, "offset": offset}); err != nil {
		return
	}
}

func queryInt(value string, fallback int) (int, error) {
	if value == "" {
		return fallback, nil
	}
	return strconv.Atoi(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}
