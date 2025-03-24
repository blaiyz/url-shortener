package main

import (
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"strings"

	"github.com/lmittmann/tint"
)

const minArgCount = 2

func main() {
	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	if len(os.Args) < minArgCount {
		slog.Error("Usage: url-shortener <URL>")
		os.Exit(1)
	}

	u := strings.TrimSpace(os.Args[1])
	parsed, err := url.Parse(u)
	if err != nil {
		slog.Error(fmt.Sprintf("Invalid URL: %s", u))
		os.Exit(1)
	}

	if !parsed.IsAbs() {
		slog.Error(fmt.Sprintf("Not an absolute URL: %s", parsed))
		os.Exit(1)
	}

	slog.Info(fmt.Sprintf("Received URL: %v", parsed))
	StartServer(parsed)
}
