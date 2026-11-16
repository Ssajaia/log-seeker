package main

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

type Store struct{ db *sql.DB }

func OpenStore(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	db.SetMaxOpenConns(1)
	if _, err := db.Exec(`PRAGMA busy_timeout = 5000; CREATE TABLE IF NOT EXISTS logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp TEXT NOT NULL,
		level TEXT NOT NULL,
		source TEXT NOT NULL,
		message TEXT NOT NULL
	); CREATE INDEX IF NOT EXISTS idx_logs_timestamp ON logs(timestamp DESC, id DESC);
	CREATE INDEX IF NOT EXISTS idx_logs_level ON logs(level);`); err != nil {
		db.Close()
		return nil, fmt.Errorf("initialize database: %w", err)
	}
	return &Store{db: db}, nil
}

func (s *Store) Close() error { return s.db.Close() }

func (s *Store) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM logs`).Scan(&count); err != nil {
		return 0, fmt.Errorf("count logs: %w", err)
	}
	return count, nil
}

func (s *Store) Insert(ctx context.Context, entry LogEntry) (int64, error) {
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now().UTC()
	}
	entry.Level = strings.ToUpper(strings.TrimSpace(entry.Level))
	result, err := s.db.ExecContext(ctx, `INSERT INTO logs(timestamp, level, source, message) VALUES(?, ?, ?, ?)`, entry.Timestamp.UTC().Format(time.RFC3339Nano), entry.Level, entry.Source, entry.Message)
	if err != nil {
		return 0, fmt.Errorf("insert log: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get inserted log id: %w", err)
	}
	return id, nil
}

func (s *Store) InsertMany(ctx context.Context, entries []LogEntry) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin log batch: %w", err)
	}
	defer tx.Rollback()

	statement, err := tx.PrepareContext(ctx, `INSERT INTO logs(timestamp, level, source, message) VALUES(?, ?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("prepare log batch: %w", err)
	}
	defer statement.Close()
	for _, entry := range entries {
		if entry.Timestamp.IsZero() {
			entry.Timestamp = time.Now().UTC()
		}
		entry.Level = strings.ToUpper(strings.TrimSpace(entry.Level))
		if _, err := statement.ExecContext(ctx, entry.Timestamp.UTC().Format(time.RFC3339Nano), entry.Level, entry.Source, entry.Message); err != nil {
			return fmt.Errorf("insert log batch: %w", err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit log batch: %w", err)
	}
	return nil
}

func (s *Store) List(ctx context.Context, filter LogFilter) ([]LogEntry, error) {
	query := `SELECT id, timestamp, level, source, message FROM logs`
	var clauses []string
	var args []any
	if filter.Level != "" {
		clauses = append(clauses, `level = ?`)
		args = append(args, strings.ToUpper(filter.Level))
	}
	var matcher *regexp.Regexp
	if filter.Regex != "" {
		var err error
		matcher, err = regexp.Compile(filter.Regex)
		if err != nil {
			return nil, fmt.Errorf("compile log regex: %w", err)
		}
	}
	if len(clauses) > 0 {
		query += ` WHERE ` + strings.Join(clauses, ` AND `)
	}
	query += ` ORDER BY timestamp DESC, id DESC`
	if matcher == nil {
		query += ` LIMIT ? OFFSET ?`
		args = append(args, filter.Limit, filter.Offset)
	}
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query logs: %w", err)
	}
	defer rows.Close()
	entries := make([]LogEntry, 0)
	matched := 0
	for rows.Next() {
		var entry LogEntry
		var timestamp string
		if err := rows.Scan(&entry.ID, &timestamp, &entry.Level, &entry.Source, &entry.Message); err != nil {
			return nil, fmt.Errorf("read log: %w", err)
		}
		entry.Timestamp, err = time.Parse(time.RFC3339Nano, timestamp)
		if err != nil {
			return nil, fmt.Errorf("parse stored timestamp: %w", err)
		}
		if matcher != nil {
			if !matcher.MatchString(entry.Message) {
				continue
			}
			if matched < filter.Offset {
				matched++
				continue
			}
			matched++
		}
		entries = append(entries, entry)
		if matcher != nil && len(entries) >= filter.Limit {
			break
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate logs: %w", err)
	}
	return entries, nil
}
