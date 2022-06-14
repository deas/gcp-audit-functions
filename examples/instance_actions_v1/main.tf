locals {
  name       = "instance-actions"
  start_name = "start"
  stop_name  = "stop"
  // TODO: Some action bits should probably go down to functions module
  action = {
    "stop"  = <<-EOF
        {"action": "stop",
         "search": {
           "scope": "${var.search_scope}/${var.project_id}",
           "query": "${var.action.stop.query}",
           "assetTypes": ["compute.googleapis.com/Instance"]
         }
        }
      EOF
    "start" = <<-EOF
        {"action": "start",
         "search": {
           "scope": "${var.search_scope}/${var.project_id}",
           "query": "${var.action.start.query}",
           "assetTypes": ["compute.googleapis.com/Instance"]
         }
        }
      EOF
  }
}

resource "random_pet" "main" {
  separator = "-"
}

# TODO: A bit hacky? Open for ideas.
module "function" {
  source = "../.."
}


# Shorter if you just want one (-> message) scheduled function invocation
/*
module "start_instances_scheduled" {
  source                    = "terraform-google-modules/scheduled-function/google"
  version                   = "2.4.0"
  project_id                = var.project_id
  job_name                  = "${local.name}-${random_pet.main.id}"
  job_description           = "Scheduled Start of GCE Instances based on Asset Search"
  job_schedule              = var.schedule
  function_entry_point      = module.function.entry_points_v1["instance_actions"]
  function_source_directory = module.function.path
  # function_source_dependent_files = [local_file.package] #, local_file.shadow]
  function_name                = "${local.name}-${random_pet.main.id}"
  function_description         = "Start Instances based on Asset Search"
  region                       = var.region
  topic_name                   = "${local.name}-${random_pet.main.id}"
  function_runtime             = module.function.runtime
  function_available_memory_mb = "128"
  #function_secret_environment_variables = [
  #  {
  #    key         = "SECRET_KEY"
  #    project_id  = var.secret_project_id
  #    secret_name = var.secret_name
  #    version     = "1"
  #  }
  #]
  #function_environment_variables = {
  #  A = "b"
  #}
  message_data = base64encode(var.search)
}
*/

# Only for reference - soft delete of roles does not play nice with terraform
/*
resource "google_organization_iam_custom_role" "instance_actions" {
  role_id     = "ComputeInstancesActions"
  title       = "Compute Instances Actions"
  description = "Compute Instances Actions"
  org_id      = "746882147037"
  permissions = [
    "compute.instances.start",
    "compute.instances.stop",
    "cloudasset.assets.searchAllResources"
  ]
}

resource "google_organization_iam_binding" "appspot" {
  org_id  = "746882147037"
  role    = google_organization_iam_custom_role.instance_actions.id
  members = [
    "serviceAccount:${var.project_id}@appspot.gserviceaccount.com"
  ]
}
*/

resource "google_pubsub_topic" "main" {
  name                       = "${local.name}-${random_pet.main.id}"
  project                    = var.project_id
  message_retention_duration = "86600s"
}

resource "google_cloud_scheduler_job" "start" {
  name        = "${local.start_name}-${random_pet.main.id}"
  project     = var.project_id
  region      = var.region
  description = "Start VM instances"
  schedule    = var.action[local.start_name]["schedule"]
  time_zone   = var.time_zone

  pubsub_target {
    topic_name = "projects/${var.project_id}/topics/${google_pubsub_topic.main.name}"
    data       = base64encode(local.action[local.start_name])
  }
}

resource "google_cloud_scheduler_job" "stop" {
  name        = "${local.stop_name}-${random_pet.main.id}"
  project     = var.project_id
  region      = var.region
  description = "Stop VM instances"
  schedule    = var.action[local.stop_name]["schedule"]
  time_zone   = var.time_zone

  pubsub_target {
    topic_name = "projects/${var.project_id}/topics/${google_pubsub_topic.main.name}"
    data       = base64encode(local.action[local.stop_name])
  }
}

module "instance_actions" {
  source      = "terraform-google-modules/event-function/google"
  version     = "2.3.0"
  description = "Start/Stop Instance actions based on Asset Search"
  entry_point = module.function.entry_points_v1["instance_actions"]
  #environment_variables = {
  #  LABEL_KEY = "principal-email"
  #}
  event_trigger = {
    event_type = "google.pubsub.topic.publish"
    resource   = google_pubsub_topic.main.id
  }
  event_trigger_failure_policy_retry = false
  # module.event_log_entry.function_event_trigger
  name                           = "${local.name}-${random_pet.main.id}"
  project_id                     = var.project_id
  region                         = var.region
  source_directory               = module.function.path
  files_to_exclude_in_source_dir = module.function.excludes
  available_memory_mb            = "128"
  runtime                        = module.function.runtime
  service_account_email          = var.service_account_email
}