locals {
  name = "start-instances"
}

resource "random_pet" "main" {
  separator = "-"
}

# TODO: A bit hacky? Open for ideas.
module "function" {
  source = "../.."
}


/*
module "start_instances_scheduled" {
  source                    = "terraform-google-modules/scheduled-function/google"
  version                   = "2.4.0"
  project_id                = var.project_id
  job_name                  = "${local.name}-${random_pet.main.id}"
  job_description           = "Scheduled Start of GCE Instances based on Asset Search"
  job_schedule              = var.schedule
  function_entry_point      = module.function.entry_points_v1["start_instances"]
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

resource "google_pubsub_topic" "main" {
  name                       = "${local.name}-${random_pet.main.id}"
  project                    = var.project_id
  message_retention_duration = "86600s"
}

module "start_instances" {
  source      = "terraform-google-modules/event-function/google"
  version     = "2.3.0"
  description = "Start Instances based on Asset Search"
  entry_point = module.function.entry_points_v1["start_instances"]
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


/* Hopefully no longer needed
resource "null_resource" "wait_for_function" {
  provisioner "local-exec" {
    command = "sleep 60"
  }

  depends_on = [module.start_instances]
}
*/