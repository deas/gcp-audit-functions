locals {
  name = "audit-label"
}
resource "random_pet" "main" {
  separator = "-"
}

# TODO: A bit hacky? Open for ideas.
module "function" {
  source = "../.."
}

# https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/eventarc_trigger#destination
/*
gcloud alpha functions deploy gce-vm-labeler \
  --v2 \
  --runtime nodejs14 \
  --entry-point labelVmCreation \
  --source . \
  --trigger-event-filters="type=google.cloud.audit.log.v1.written,serviceName=compute.googleapis.com,methodName=beta.compute.instances.insert" \
  --region us-west1 \
  --trigger-location us-central1
*/

resource "google_eventarc_trigger" "main" {
  name     = "${local.name}-${random_pet.main.id}"
  location = var.region
  matching_criteria {
    attribute = "type"
    value     = "google.cloud.audit.log.v1.written"
  }
  destination {
    # (Optional) [WARNING] Configuring a Cloud Function in Trigger is not supported as of today.
    # The Cloud Function resource name. Format: projects/{project}/locations/{location}/functions/{function}
    cloud_function = "projects/${var.project_id}/locations/${var.region}/functions/audit-label-${random_pet.main.id}"
  }
}

module "audit_label" {
  source      = "../../modules/functions-v2"
  description = "Labels resource with owner information."
  entry_point = module.function.v2_entry_point
  #environment_variables = {
  #  LABEL_KEY = "principal-email"
  #}
  event_trigger                  = google_eventarc_trigger.main.name
  name                           = "${local.name}-${random_pet.main.id}"
  project_id                     = var.project_id
  region                         = var.region
  source_directory               = module.function.path
  files_to_exclude_in_source_dir = module.function.excludes
  available_memory_mb            = "128"
  runtime                        = module.function.runtime
}

resource "null_resource" "wait_for_function" {
  provisioner "local-exec" {
    command = "sleep 60"
  }

  depends_on = [module.audit_label]
}