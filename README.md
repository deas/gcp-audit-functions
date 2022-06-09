# GCP Cloud Functions driven by Schedule and Audit Events

This project aims to provide generic schedule and audit event driven Cloud Functions.

Functionality currently covers:

- Labeling GCE instances on creation
- Hardening the Compute Default account (revoking `role/editor`)
- Start and Stop GCE instances based on Asset Search

More hopefully coming soon.

Additionally, we aim at decent support for the larger product lifecyle with an emphasis on a DevOps experience including short cycle times. We leverage Cloud Foundation Toolkit, Cloud Functions Framework, GitHub Actions and Terraform. We cover unit- and integration testing. We stripped dependencies where reasonable and extended where we wanted to go further or connected the dots.

The v1 versions leverage PubSub Log Sinks, 🧪 v2 🥼 is based on EventArc/CloudEvents.

## Usage
Sample Cloud Function and VM deployments designed to play together are provided in the `examples` folder. Unless explicitly disabled, they are also used by the integration tests.

You may want to
```shell
export GOOGLE_IMPERSONATE_SERVICE_ACCOUNT=your-sa@your-prj-id.iam.gserviceaccount.com
```
to get proper access when trying them out.

<!-- BEGINNING OF PRE-COMMIT-TERRAFORM DOCS HOOK -->
### Inputs

No input.

### Outputs

| Name | Description |
|------|-------------|
| entry\_points\_v1 | The v1 function entry points provided by this module |
| excludes | Files we want to exlude |
| path | The path to the function source |
| runtime | The runtime |
| v1\_entry\_point | The v1 legacy label function entry point |
| v2\_entry\_point | The v2 legacy label function entry point |

<!-- END OF PRE-COMMIT-TERRAFORM DOCS HOOK -->

## Development
There are various `Makefile` targets providing entrypoints for CI and steps you might want to do during development.

Cloud Function implementations are currently Go based and we use [Functions Framework for Go](https://github.com/GoogleCloudPlatform/functions-framework-go) during development.

Start local service
```shell
# export FUNCTION_TARGET=LabelPubSub # Not needed atm
# export GCP_HOUSEKEEPER_READ_ONLY=1 # If you want read only access to GCP 
export GCP_HOUSEKEEPER_FUNCTION=StartPubSub # Framework workaround atm
make serve
```

Send PubSub payload to local Label Function
```shell
message=test/audit-compute-instance-create.json
endpoint=http://localhost:8080 # Issue with framework : Only one endpoint per process

cat <<EOF | curl -d @- -X POST -H "Content-Type: application/json" "${endpoint}" 
{
  "data": "$(base64 -w 0 $message)",
  "messageId": "3756413890745862",
  "publishTime": "2022-01-03T10:26:11.735Z"
}
EOF
```

Send CloudEvent payload to local Label Function
```shell
message=test/audit-compute-instance-create.json
endpoint=http://localhost:8080
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

Call Start Function with organization scope on GCP directly.

```shell
# Grant Function Service Account permission to search, start and stop.
ORG_ID=your-org-id
PROJECT_ID=your-project-id

gcloud iam roles create --organization ${ORG_ID} ComputeInstancesLifeCycle --permissions=compute.instances.start,compute.instances.stop,cloudasset.assets.searchAllResources

gcloud organizations add-iam-policy-binding $ORG_ID --member="serviceAccount:${PROJECT_ID}@appspot.gserviceaccount.com" --role="organizations/$ORG_ID/roles/ComputeInstancesLifeCycle"

# Find deployed function
function=$(gcloud --project ${PROJECT_ID} functions list | grep ^start.instances | cut -d " " -f 1)
gcloud --project=${PROJECT_ID} functions call ${function} --region europe-west2 --data='{"data":"'$(echo '
{
  "scope": "organizations/'$ORG_ID'",
  "query": "labels.start_daily:true AND state:TERMINATED",
  "assetTypes": ["compute.googleapis.com/Instance"]
}' | base64 -w 0)'"}'
# gcloud --project ${PROJECT_ID} pubsub topics ${function} --message '{ "fix": "me" }'
```

Alternatively, you can use `"scope": "projects/your-project-id"` or on folder level. 

Call local Start/Stop Function
```shell
endpoint=http://localhost:8080
scope=organizations/$ORG_ID

echo '{"data": "'$(echo '
{
    "scope": "'$scope'",
    "query": "labels.start_daily:true AND state:TERMINATED",
    "assetTypes": ["compute.googleapis.com/Instance"]
}' | base64 -w 0)'"}' | curl -d @- -X POST -H "Content-Type: application/json" "${endpoint}"

echo '{"data": "'$(echo '
{
    "scope": "'$scope'",
    "query": "labels.stop_daily:true AND state:RUNNING",
    "assetTypes": ["compute.googleapis.com/Instance"]
}' | base64 -w 0)'"}' | curl -d @- -X POST -H "Content-Type: application/json" "${endpoint}"

```

Read Audit Logs from StackDriver
```shell
gcloud logging read 'protoPayload.@type="type.googleapis.com/google.cloud.audit.AuditLog"' --freshness=1h --project ${PROJECT_ID} --format json
```

There is a generic `main` entrypoint for the Go implementations which allows one to call functionality straight from the command line. Try
```shell
cd fn
go run main/main.go --help
```
to see what is available.
# Known Issues
- Go workspaces are recommended for best DX with `gopls` : [x/tools/gopls: support multi-module workspaces #32394](https://github.com/golang/go/issues/32394) / [Setting up your workspace](https://github.com/golang/tools/blob/master/gopls/doc/workspace.md#go-workspaces-go-118)
- [Serving multiple functions locally from a single server instance #109](https://github.com/GoogleCloudPlatform/functions-framework-go/issues/109)
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
