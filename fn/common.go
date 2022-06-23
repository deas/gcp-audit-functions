package function

import (
	"context"
	"fmt"
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
	_, readOnly = os.LookupEnv(fmt.Sprintf("%s_%s", EnvPrefix, "READ_ONLY"))
}

func NewLogger() slog.Logger {
	// https://cloud.google.com/functions/docs/configuring/env-var#newer_runtimes
	if _, exists := os.LookupEnv("FUNCTION_SIGNATURE_TYPE"); exists {
		return slog.Make(slogstackdriver.Sink(os.Stdout))
	} else {
		return slog.Make(sloghuman.Sink(os.Stdout))
	}
}

// PubSubMessage is the payload of a Pub/Sub event - Save a module for the moment
// See the documentation for more details:
// https://cloud.google.com/pubsub/docs/reference/rest/v1/PubsubMessage
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// Cloud Scheduler -> PubSub -> Cloud Event Function (v2) Glue
type PubSubEventData struct {
	Message PubSubMessage `json:"message"`
}
