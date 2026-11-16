package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	if len(os.Args) > 2 || (len(os.Args) == 2 && os.Args[1] != "mock-logs") {
		log.Fatalf("usage: %s [mock-logs]", os.Args[0])
	}
	dbPath := envOr("LOGSEEKER_DB", "logseeker.db")
	if len(os.Args) > 1 && os.Args[1] == "mock-logs" {
		if err := seedMockLogs(dbPath); err != nil {
			log.Fatal(err)
		}
		return
	}
	if err := runServer(envOr("LOGSEEKER_ADDR", ":8080"), dbPath); err != nil {
		log.Fatal(err)
	}
}

func runServer(addr, dbPath string) error {
	store, err := OpenStore(dbPath)
	if err != nil {
		return err
	}
	defer store.Close()

	mux := http.NewServeMux()
	mux.Handle("/api/logs", NewAPI(store))
	mux.Handle("/", http.FileServer(http.Dir("frontend")))
	log.Printf("LogSeeker listening on %s (database %s)", addr, dbPath)
	server := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func seedMockLogs(dbPath string) error {
	store, err := OpenStore(dbPath)
	if err != nil {
		return err
	}
	defer store.Close()

	ctx := context.Background()
	count, err := store.Count(ctx)
	if err != nil {
		return err
	}
	if count > 0 {
		log.Printf("Database %s already has %d log(s); skipped mock data", dbPath, count)
		return nil
	}
	if err := store.InsertMany(ctx, mockLogEntries()); err != nil {
		return err
	}
	log.Printf("Inserted 30 mock logs into %s", dbPath)
	return nil
}

func envOr(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
