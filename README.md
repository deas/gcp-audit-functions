# GCP Cloud Functions triggered by Audit Events

This project aims to provide audit event triggered Cloud Functions. The functionality is provided as a `terraform` module.

The v1 versions leverage PubSub Log Sinks, 🧪 v2 🥼 is based on EventArc/CloudEvents.

We try hard to build upon people committed to maintainance (mostly CFT by Google). We stripped where possible, extended where we wanted to go further or connected the dots.

## Usage
Sample Cloud Function and VM deployments designed to play together are provided in the `examples` folder. Unless explicitly disabled, they are also used by the integration tests.

You may want to
```shell
export GOOGLE_IMPERSONATE_SERVICE_ACCOUNT=your-sa@your-prj-id.iam.gserviceaccount.com
```
to get proper access when trying them out.

### Outputs

| Name | Description |
|------|-------------|
| <a name="output_excludes"></a> [excludes](#output\_excludes) | Files we want to exlude |
| <a name="output_path"></a> [path](#output\_path) | The path to the function source |
| <a name="output_runtime"></a> [runtime](#output\_runtime) | The runtime |
| <a name="output_v1_entry_point"></a> [v1\_entry\_point](#output\_v1\_entry\_point) | The v1 function entry point |
| <a name="output_v2_entry_point"></a> [v2\_entry\_point](#output\_v2\_entry\_point) | The v2 function entry point |

## Development
There are various `Makefile` targets providing entrypoints for CI and steps you might want to do during development.

Cloud Function implementations are currently Go based and we use [Functions Framework for Go](https://github.com/GoogleCloudPlatform/functions-framework-go) during development.

Start local service
```shell
# export FUNCTION_TARGET=LabelPubSub # Not needed atm
# export GCP_AUDIT_LABEL_READ_ONLY=1 # If you want read only access to GCP 
make serve
```

Send PubSub payload to local label Function
```shell
message=test/audit-compute-instance-create.json
endpoint=http://localhost:8080/label-pubsub

cat <<EOF | curl -d @- -X POST -H "Content-Type: application/json" "${endpoint}" 
{
  "data": "$(base64 -w 0 $message)",
  "messageId": "3756413890745862",
  "publishTime": "2022-01-03T10:26:11.735Z"
}
EOF
```

Send CloudEvent payload to local Function
```shell
message=test/audit-compute-instance-create.json
endpoint=http://localhost:8080/label-event
cat <<EOF | curl -d @- -X POST -H "Content-Type: application/cloudevents+json" "${endpoint}" \
{
	"specversion" : "1.0",
	"type" : "example.com.cloud.event",
	"source" : "https://example.com/cloudevents/pull",
	"subject" : "123",
	"id" : "A234-1234-1234",
	"time" : "2018-04-05T17:31:00Z",
	"data" : $(cat $message)
}
EOF
```

Send PubSub payload to Cloud Function via Topic.

```shell
PROJECT_ID=your-prj-id
function=$(gcloud --project ${PROJECT_ID} functions list | grep ^audit-label | cut -d " " -f 1)
gcloud --project ${PROJECT_ID} functions call ${function} --data='{"message": "Hello World!"}'
gcloud --project ${PROJECT_ID} pubsub topics audit-label --message '{ "fix": "me" }'
```

Read Audit Logs from StackDriver
```shell
gcloud logging read 'protoPayload.@type="type.googleapis.com/google.cloud.audit.AuditLog"' --freshness=1h --project ${PROJECT_ID} --format json
```
# Known Issues
- Go workspaces are recommended for best DX with `gopls` : [x/tools/gopls: support multi-module workspaces #32394](https://github.com/golang/go/issues/32394) / [Setting up your workspace](https://github.com/golang/tools/blob/master/gopls/doc/workspace.md#go-workspaces-go-118)
# References
- [Functions Framework for Go](https://github.com/GoogleCloudPlatform/functions-framework-go)
- [Go SDK for CloudEvents](https://github.com/cloudevents/sdk-go).
- [Structuring source code](https://cloud.google.com/functions/docs/writing/#structuring_source_code)
- [CloudEvents Spec](https://cloudevents.io/)
- [Using Cloud Audit Logs](https://cloud.google.com/eventarc/docs/reference/supported-events#using-cloud-audit-logs)
- [Compute Engine VM Labeler - Cloud Functions v2](https://github.com/GoogleCloudPlatform/eventarc-samples/tree/main/gce-vm-labeler/gcf)
- [Cross Project Eventing](https://github.com/GoogleCloudPlatform/eventarc-samples/tree/main/cross-project-eventing)
- [Calling Cloud Functions (v1)](https://cloud.google.com/functions/docs/calling)
- [`golang-samples/functions/functionsv2`](https://github.com/GoogleCloudPlatform/golang-samples/tree/main/functions/functionsv2)
- [Cloud Foundation Toolkit Project](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit)
- [IAM audit logging ](https://cloud.google.com/iam/docs/audit-logging)
- [Audit logs for service accounts](https://cloud.google.com/iam/docs/audit-logging/examples-service-accounts)
