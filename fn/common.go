package function

import (
	"context"
	"log"
	"os"

	"cdr.dev/slog"
	"cdr.dev/slog/sloggers/sloghuman"
	"cdr.dev/slog/sloggers/slogstackdriver"
)

// https://gist.github.com/salrashid123/62178224324ccbc80a358920d5281a60

var (
	readOnly  bool
	EnvPrefix = "GCP_HOUSEKEEPER"
	Logger    slog.Logger
)

func init() {
	Logger = NewLogger()
	ctx := context.Background()
	log.SetOutput(slog.Stdlib(ctx, Logger).Writer())
	log.Print("Logging initialized")
}

func NewLogger() slog.Logger {
	// https://cloud.google.com/functions/docs/configuring/env-var#newer_runtimes
	if _, exists := os.LookupEnv("FUNCTION_SIGNATURE_TYPE"); exists {
		return slog.Make(slogstackdriver.Sink(os.Stdout))
	} else {
		return slog.Make(sloghuman.Sink(os.Stdout))
	}
}
