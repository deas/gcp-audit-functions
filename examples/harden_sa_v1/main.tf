locals {
  name = "audit-harden-sa"
}

resource "random_pet" "main" {
  separator = "-"
}

# TODO: Delegates to
# https://github.com/terraform-google-modules/terraform-google-log-export
# Event function submodules also support folders - but not organization

module "event_log_entry" {
  # https://cloud.google.com/iam/docs/audit-logging
  # operation.first=true
  source     = "../../modules/event-organization-log-entry"
  filter     = "protoPayload.@type=type.googleapis.com/google.cloud.audit.AuditLog AND (protoPayload.methodName=google.iam.admin.v1.CreateServiceAccount OR protoPayload.methodName=SetIamPolicy)"
  name       = "${local.name}-${random_pet.main.id}"
  project_id = var.project_id
  org_id     = var.org_id
}

/*
module "event_log_entry" {
  source     = "terraform-google-modules/event-function/google//modules/event-project-log-entry"
  version    = "2.3.0"
  filter          = "protoPayload.@type=\"type.googleapis.com/google.cloud.audit.AuditLog\" protoPayload.methodName=google.iam.admin.v1.CreateServiceAccount"
  name       = "${local.name}-${random_pet.main.id}"
  project_id = var.project_id
}
*/

# TODO: A bit hacky? Open for ideas.
module "function" {
  source = "../.."
}

module "harden_sa" {
  source      = "terraform-google-modules/event-function/google"
  version     = "2.3.0"
  description = "Harden Default Compute Service Account Policy Binding."
  entry_point = module.function.entry_points_v1["harden_sa"]
  #environment_variables = {
  #  LABEL_KEY = "principal-email"
  #}
  event_trigger                  = module.event_log_entry.function_event_trigger
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

  depends_on = [module.harden_sa]
}