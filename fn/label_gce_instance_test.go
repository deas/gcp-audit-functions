package function

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudevents/sdk-go/v2/event"
	"google.golang.org/api/option"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

type FakeInstancesServer struct {
	http.ServeMux
	changedLabels map[string]string
	instance      *computepb.Instance
}

func NewFakeInstancesServer(i *computepb.Instance) *FakeInstancesServer {
	fake := &FakeInstancesServer{
		instance:      i,
		changedLabels: make(map[string]string),
	}
	fake.HandleFunc("/compute/v1/projects/PROJECT/zones/ZONE/instances/INSTANCE",
		func(w http.ResponseWriter, r *http.Request) {
			rbytes, err := protojson.Marshal(fake.instance)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			w.Write(rbytes)
		})
	fake.HandleFunc("/compute/v1/projects/PROJECT/zones/ZONE/instances/INSTANCE/setLabels",
		func(w http.ResponseWriter, r *http.Request) {
			var labelReq computepb.InstancesSetLabelsRequest
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			err = protojson.Unmarshal(body, &labelReq)
			// var audit_data audit.AuditLog
			// err := anypb.UnmarshalTo(e.GetProtoPayload(), &audit_data, proto.UnmarshalOptions{})

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			fake.changedLabels = labelReq.GetLabels()
			response, _ := protojson.Marshal(&computepb.Operation{})
			w.Write(response)
			w.WriteHeader(http.StatusOK)
		})
	fake.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})
	return fake
}

func TestLabelGceInstance(t *testing.T) {
	tests := []struct {
		name       string
		payload    *AuditLogProtoPayload
		instance   *computepb.Instance
		newCreator string
	}{
		{
			name: "no label",
			payload: &AuditLogProtoPayload{
				MethodName:   "beta.compute.instances.insert",
				ResourceName: "projects/abc-0815/zones/europe-west3-b/instances/test",
				ServiceName:  "compute.googleapis.com",
				AuthenticationInfo: map[string]interface{}{
					"principalEmail": "user@example.com",
				},
			},
			instance: &computepb.Instance{
				Hostname:         new(string),
				Id:               new(uint64),
				Kind:             new(string),
				LabelFingerprint: new(string),
				Labels:           map[string]string{},
				Name:             new(string),
				SelfLink:         new(string),
				Zone:             new(string),
			},
			newCreator: "user_example_com",
		},
		{
			name: "existing-creator",
			payload: &AuditLogProtoPayload{
				MethodName:   "beta.compute.instances.insert",
				ResourceName: "projects/abc-0815/zones/europe-west3-b/instances/test",
				ServiceName:  "compute.googleapis.com",
				AuthenticationInfo: map[string]interface{}{
					"principalEmail": "user@example.com",
				},
			},
			instance: &computepb.Instance{
				Hostname:         new(string),
				Id:               new(uint64),
				Kind:             new(string),
				LabelFingerprint: new(string),
				Labels: map[string]string{
					"creator": "existing_creator",
				},
				Name:     new(string),
				SelfLink: new(string),
				Zone:     new(string),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fake := NewFakeInstancesServer(tt.instance)
			fakeserver := httptest.NewServer(fake)
			defer fakeserver.Close()

			var err error
			Opts = []option.ClientOption{
				option.WithEndpoint(fakeserver.URL),
				option.WithoutAuthentication()}
			/* client, err = compute.NewInstancesRESTClient(context.Background(), option.WithEndpoint(fakeserver.URL), option.WithoutAuthentication())
			if err != nil {
				t.Fatalf("Failed to create mock client: %s", err)
			} */
			e := event.New()
			e.SetSubject("compute.googleapis.com/projects/PROJECT/zones/ZONE/instances/INSTANCE")
			e.SetType("google.cloud.audit.log.v1.written")
			auditlog := &AuditLogEntry{
				ProtoPayload: tt.payload,
			}
			eventdata, err := json.Marshal(auditlog)
			if err != nil {
				t.Fatalf("Failed to marshal json for test %s: %s", tt.name, err)
			}
			e.SetDataContentType("application/json")
			e.SetData("application/json", eventdata)

			if err := label(context.Background(), *auditlog.ProtoPayload, e.Subject()); err != nil {
				t.Fatalf("LabelEventInstance(%s): unexpected error %s", tt.name, err)
			}
			// check that we updated creator label if expected.
			if tt.newCreator != "" {
				newvalue, ok := fake.changedLabels["creator"]
				if !ok || newvalue != tt.newCreator {
					t.Fatalf("LabelEventInstance(%s): incorrect creator label: got %s, want %s", tt.name, newvalue, tt.newCreator)
				}
			}

		})
	}
}
