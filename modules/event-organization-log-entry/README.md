# Event Organization Log Entry

This submodule configures a organisation-level Stackdriver Logging export to
act as an event which will trigger a Cloud Functions function.

The export uses a provided filter to identify events of interest and
publishes them to a dedicated Pub/Sub topic. The target function
must be configured to subscribe to the topic in order to process each
export event.

Disclaimer: Code mostly borrowed from `terraform-google-event-function`.

## Usage

...
<!-- BEGINNING OF PRE-COMMIT-TERRAFORM DOCS HOOK -->
## Inputs
TODO: Rebuild

## Outputs

TODO: Rebuild

<!-- END OF PRE-COMMIT-TERRAFORM DOCS HOOK -->

## Requirements

The following sections describe the requirements which must be met in
order to invoke this module.

### Software Dependencies

The following software dependencies must be installed on the system
from which this module will be invoked:

- [Terraform][terraform-site] v0.12
- [Terraform Provider for Google Cloud Platform][terraform-provider-gcp-site] v2.5

### IAM Roles

The Service Account which will be used to invoke this module must have
the following IAM roles:

- Logs Configuration Writer: `roles/logging.configWriter`
- Pub/Sub Admin: `roles/pubsub.admin`
- Service Account User: `roles/iam.serviceAccountUser`

### APIs

The project against which this module will be invoked must have the
following APIs enabled:

- Cloud Pub/Sub API: `pubsub.googleapis.com`
- Stackdriver Logging API: `logging.googleapis.com`