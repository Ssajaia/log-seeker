package main

import "time"

// mockLogEntries provides a small, clearly labeled dataset for local demos.
func mockLogEntries() []LogEntry {
	base := time.Now().UTC().Truncate(time.Second)
	rows := []struct{ level, source, message string }{
		{"INFO", "api-gateway", "Mock log 01: service started"},
		{"INFO", "auth-service", "Mock log 02: user session created"},
		{"DEBUG", "worker-1", "Mock log 03: polling job queue"},
		{"WARN", "api-gateway", "Mock log 04: request took longer than expected"},
		{"INFO", "database", "Mock log 05: connection pool ready"},
		{"ERROR", "worker-2", "Mock log 06: failed to process sample task"},
		{"INFO", "scheduler", "Mock log 07: scheduled task completed"},
		{"DEBUG", "worker-1", "Mock log 08: retry counter is 1"},
		{"WARN", "database", "Mock log 09: connection pool is nearly full"},
		{"INFO", "auth-service", "Mock log 10: token validation succeeded"},
		{"ERROR", "api-gateway", "Mock log 11: upstream returned status 503"},
		{"INFO", "worker-2", "Mock log 12: sample task accepted"},
		{"DEBUG", "scheduler", "Mock log 13: next run calculated"},
		{"WARN", "auth-service", "Mock log 14: login attempt needs review"},
		{"INFO", "database", "Mock log 15: query completed in 18ms"},
		{"ERROR", "worker-1", "Mock log 16: temporary file was not found"},
		{"INFO", "api-gateway", "Mock log 17: request completed with status 200"},
		{"DEBUG", "worker-2", "Mock log 18: task payload decoded"},
		{"WARN", "scheduler", "Mock log 19: skipped overlapping run"},
		{"INFO", "auth-service", "Mock log 20: user session closed"},
		{"ERROR", "database", "Mock log 21: sample transaction rolled back"},
		{"INFO", "worker-1", "Mock log 22: retry completed successfully"},
		{"DEBUG", "api-gateway", "Mock log 23: response headers prepared"},
		{"WARN", "worker-2", "Mock log 24: task is taking extra time"},
		{"INFO", "scheduler", "Mock log 25: scheduled task started"},
		{"ERROR", "auth-service", "Mock log 26: sample token was rejected"},
		{"INFO", "database", "Mock log 27: health check passed"},
		{"DEBUG", "worker-1", "Mock log 28: no additional jobs available"},
		{"WARN", "api-gateway", "Mock log 29: response used fallback route"},
		{"INFO", "logseeker", "Mock log 30: mock data batch is ready"},
	}

	entries := make([]LogEntry, 0, len(rows))
	for i, row := range rows {
		entries = append(entries, LogEntry{
			Timestamp: base.Add(-time.Duration(len(rows)-1-i) * time.Minute),
			Level:     row.level,
			Source:    row.source,
			Message:   row.message,
		})
	}
	return entries
}
