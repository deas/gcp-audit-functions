package function

import (
	"fmt"
	"os"

	"cdr.dev/slog"
	"cdr.dev/slog/sloggers/sloghuman"
	"cdr.dev/slog/sloggers/slogstackdriver"
	"google.golang.org/api/option"
)

var (
	readOnly  bool
	EnvPrefix = "GCP_AUDIT_LABEL"
	log       slog.Logger
)

// https://gist.github.com/salrashid123/62178224324ccbc80a358920d5281a60

// AuditLogEntry represents a LogEntry as described at
// https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry
type AuditLogEntry struct {
	ProtoPayload *AuditLogProtoPayload `json:"protoPayload"`
	// ProtoPayload *audit.AuditLog `json:"protoPayload"`
}

type AuditLogProtoPayload struct {
	// protobuf version not compatible it seems - need to check what exactly sink drops on the topic
	// MethodName string `protobuf:"bytes,8,opt,name=method_name,json=methodName,proto3" json:"method_name,omitempty"`
	ServiceName        string                 `json:"serviceName"`
	MethodName         string                 `json:"methodName"`
	ResourceName       string                 `json:"resourceName"`
	AuthenticationInfo map[string]interface{} `json:"authenticationInfo"`
}

// PubSubMessage is the payload of a Pub/Sub event - Save a module for the moment
// See the documentation for more details:
// https://cloud.google.com/pubsub/docs/reference/rest/v1/PubsubMessage
type PubSubMessage struct {
	Data []byte `json:"data"`
}

func NewLogger() slog.Logger {
	if os.Getenv(fmt.Sprintf("%s_%s", EnvPrefix, "LOGGER")) == "human" {
		return slog.Make(sloghuman.Sink(os.Stdout))
	} else {
		return slog.Make(slogstackdriver.Sink(os.Stdout))
	}
}

func NewOpts() []option.ClientOption {
	var opts []option.ClientOption
	// Don't up in Pipeline : google: could not find default credentials.
	if _, noAuth := os.LookupEnv(fmt.Sprintf("%s_%s", EnvPrefix, "NO_AUTH")); noAuth {
		opts = []option.ClientOption{
			option.WithoutAuthentication(),
		}
	} else {
		opts = []option.ClientOption{}
	}
	return opts
}
