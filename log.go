package main

import "time"

// LogEntry is the JSON representation shared by the API and frontend.
type LogEntry struct {
	ID        int64     `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Source    string    `json:"source"`
	Message   string    `json:"message"`
}

type LogFilter struct {
	Level  string
	Regex  string
	Limit  int
	Offset int
}
