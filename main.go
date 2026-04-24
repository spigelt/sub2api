package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/sub2api/sub2api/handler"
)

const (
	defaultPort    = 8080
	defaultHost    = "127.0.0.1" // changed from 0.0.0.0 to localhost-only for personal use
	defaultBaseURL = ""
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	var (
		port    int
		host    string
		baseURL string
		showVer bool
	)

	flag.IntVar(&port, "port", getEnvInt("PORT", defaultPort), "Port to listen on")
	flag.StringVar(&host, "host", getEnv("HOST", defaultHost), "Host to bind to")
	flag.StringVar(&baseURL, "base-url", getEnv("BASE_URL", defaultBaseURL), "Base URL for the API (e.g. https://example.com)")
	flag.BoolVar(&showVer, "version", false, "Print version information and exit")
	flag.Parse()

	if showVer {
		fmt.Printf("sub2api %s (commit: %s, built: %s)\n", version, commit, date)
		os.Exit(0)
	}

	log.Printf("Starting sub2api %s", version)
	log.Printf("Listening on %s:%d", host, port)

	router := handler.NewRouter(handler.Config{
		BaseURL: baseURL,
	})

	addr := fmt.Sprintf("%s:%d", host, port)
	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// getEnv returns the value of the environment variable named by the key,
// or fallback if the variable is not set.
func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

// getEnvInt returns the integer value of the environment variable named by
// the key, or fallback if the variable is not set or cannot be parsed.
func getEnvInt(key string, fallback int) int {
	if val, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return fallback
}
