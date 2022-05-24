package function

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"cdr.dev/slog"
	"cdr.dev/slog/sloggers/sloghuman"
	"cdr.dev/slog/sloggers/slogstackdriver"

	compute "cloud.google.com/go/compute/apiv1"
	"github.com/cloudevents/sdk-go/v2/event"
	"google.golang.org/api/option"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"

	// https://gist.github.com/salrashid123/62178224324ccbc80a358920d5281a60
	"google.golang.org/protobuf/proto"
)

var (
	client    *compute.InstancesClient
	readOnly  bool
	EnvPrefix = "GCP_AUDIT_LABEL"
	log       slog.Logger
)

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

func init() {
	log = NewLogger()
	var err error
	client, err = compute.NewInstancesRESTClient(context.Background(), NewOpts()...)
	if err != nil {
		log.Fatal(context.Background(), fmt.Sprintf("Failed to create instances client : %v", err))
	}
	_, readOnly = os.LookupEnv(fmt.Sprintf("%s_%s", EnvPrefix, "READ_ONLY"))
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

// PubSubMessage is the payload of a Pub/Sub event - Save a module for the moment
// See the documentation for more details:
// https://cloud.google.com/pubsub/docs/reference/rest/v1/PubsubMessage
type PubSubMessage struct {
	Data []byte `json:"data"`
}

func LabelPubSub(ctx context.Context, m PubSubMessage /* pubsub.Message*/) error {
	log.Info(ctx, "Got PubSub message") // , m.ID /*string(m.Data)*/)) // Automatically decoded from base64
	logentry := &AuditLogEntry{}
	// logentry := &audit.AuditLog{}
	// var auditLogEntry AuditLogEntry
	err := json.Unmarshal(m.Data, &logentry)
	if err != nil {
		log.Info(ctx, fmt.Sprintf("Error: could not unmarshall to audit log %v\n", err))
	}
	label(ctx, *logentry.ProtoPayload, fmt.Sprintf("%s/%s", logentry.ProtoPayload.ServiceName, logentry.ProtoPayload.ResourceName))
	return nil
}

func LabelEvent(ctx context.Context, ev event.Event) error {
	log.Info(ctx, fmt.Sprintf("Got CloudEvent %s with data %s", ev.ID(), string(ev.Data())))
	logentry := &AuditLogEntry{}
	if err := ev.DataAs(logentry); err != nil {
		return fmt.Errorf("Error parsing event payload : %w", err)
	}
	return label(ctx, *logentry.ProtoPayload, ev.Subject())
}

// Receives GCE instance creation Audit Logs, and adds a `creator` label to the instance.
func label(ctx context.Context, payload AuditLogProtoPayload, subject string) error {
	log.Info(ctx, fmt.Sprintf("Got MethodName %s", payload.MethodName))

	creator, ok := payload.AuthenticationInfo["principalEmail"]
	if !ok {
		err := fmt.Errorf("principalEmail not found in cloud event payload: %v", payload)
		log.Info(ctx, fmt.Sprintf("Creator email not found: %s", err))
		return err
	}

	// Get relevant VM instance details from the event's `subject` property
	// Subject format:
	// compute.googleapis.com/projects/<PROJECT>/zones/<ZONE>/instances/<INSTANCE>
	paths := strings.Split(subject, "/")
	if len(paths) < 6 {
		return fmt.Errorf("invalid event subject: %s", subject)
	}
	project := paths[2]
	zone := paths[4]
	instance := paths[6]

	// Sanitize the `creator` label value to match GCE label requirements
	// See https://cloud.google.com/compute/docs/labeling-resources#requirements
	labelSanitizer := regexp.MustCompile("[^a-z0-9_-]+")
	creatorstring := labelSanitizer.ReplaceAllString(strings.ToLower(creator.(string)), "_")
	// creatorstring := labelSanitizer.ReplaceAllString(strings.ToLower(creator), "_")

	// Get the newly-created VM instance's label fingerprint
	// This is a requirement of the Compute Engine API and avoids duplicate labels
	inst, err := client.Get(ctx, &computepb.GetInstanceRequest{
		Project:  project,
		Zone:     zone,
		Instance: instance,
	})
	if err != nil {
		return fmt.Errorf("Could not retrieve GCE instance: %s", err)
	}
	if v, ok := inst.Labels["creator"]; ok {
		// Instance already has a creator label.
		log.Info(ctx, fmt.Sprintf("Instance %s already labeled with creator: %s", instance, v))
		return nil
	}

	if !readOnly {
		// Add the creator label to the instance
		op, err := client.SetLabels(ctx, &computepb.SetLabelsInstanceRequest{
			Project:  project,
			Zone:     zone,
			Instance: instance,
			InstancesSetLabelsRequestResource: &computepb.InstancesSetLabelsRequest{
				LabelFingerprint: proto.String(inst.GetLabelFingerprint()),
				Labels: map[string]string{
					"creator": creatorstring,
				},
			},
		})
		if err != nil {
			return err // log.Fatalf("Could not label GCE instance: %s", err)
		}
		log.Info(ctx, fmt.Sprintf("Creator label added to %s in operation %v", instance, op))
	} else {
		log.Info(ctx, fmt.Sprintf("Creator label not added to %s - read only", instance))
	}
	return nil
}
