package main

import (
	"context"
	"testing"
	"time"
)

func TestStoreInsertListAndPagination(t *testing.T) {
	store, err := OpenStore(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	ctx := context.Background()
	times := []time.Time{time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC), time.Date(2025, 1, 3, 0, 0, 0, 0, time.UTC)}
	for i, entry := range []LogEntry{{Timestamp: times[0], Level: "info", Source: "api", Message: "started"}, {Timestamp: times[1], Level: "error", Source: "db", Message: "timeout"}, {Timestamp: times[2], Level: "ERROR", Source: "db", Message: "offline"}} {
		id, err := store.Insert(ctx, entry)
		if err != nil {
			t.Fatal(err)
		}
		if id != int64(i+1) {
			t.Fatalf("id = %d, want %d", id, i+1)
		}
	}
	logs, err := store.List(ctx, LogFilter{Limit: 1, Offset: 1})
	if err != nil {
		t.Fatal(err)
	}
	if len(logs) != 1 || logs[0].Message != "timeout" {
		t.Fatalf("unexpected page: %#v", logs)
	}
	logs, err = store.List(ctx, LogFilter{Level: "error", Limit: 10})
	if err != nil {
		t.Fatal(err)
	}
	if len(logs) != 2 || logs[0].Message != "offline" {
		t.Fatalf("unexpected severity results: %#v", logs)
	}
	logs, err = store.List(ctx, LogFilter{Regex: "timeout|offline", Limit: 1, Offset: 1})
	if err != nil {
		t.Fatal(err)
	}
	if len(logs) != 1 || logs[0].Message != "timeout" {
		t.Fatalf("regex pagination should apply after matching: %#v", logs)
	}
}
