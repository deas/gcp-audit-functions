locals {
  logging = var.log_bucket == null ? [] : [
    {
      log_bucket        = var.log_bucket
      log_object_prefix = var.log_object_prefix
    }
  ]
}

resource "null_resource" "dependent_files" {
  triggers = {
    for file in var.source_dependent_files :
    pathexpand(file.filename) => file.id
  }
}

data "null_data_source" "wait_for_files" {
  inputs = {
    # This ensures that this data resource will not be evaluated until
    # after the null_resource has been created.
    dependent_files_id = null_resource.dependent_files.id

    # This value gives us something to implicitly depend on
    # in the archive_file below.
    source_dir = pathexpand(var.source_directory)
  }
}

data "archive_file" "main" {
  type        = "zip"
  output_path = pathexpand("${var.source_directory}.zip")
  source_dir  = data.null_data_source.wait_for_files.outputs["source_dir"]
  excludes    = var.files_to_exclude_in_source_dir
}

resource "google_storage_bucket" "main" {
  count                       = var.create_bucket ? 1 : 0
  name                        = coalesce(var.bucket_name, var.name)
  force_destroy               = var.bucket_force_destroy
  location                    = var.region
  project                     = var.project_id
  storage_class               = "REGIONAL"
  labels                      = var.bucket_labels
  uniform_bucket_level_access = true

  dynamic "logging" {
    for_each = local.logging == [] ? [] : local.logging
    content {
      log_bucket        = logging.value.log_bucket
      log_object_prefix = logging.value.log_object_prefix
    }
  }

}

resource "google_storage_bucket_object" "main" {
  name                = "${data.archive_file.main.output_md5}-${basename(data.archive_file.main.output_path)}"
  bucket              = var.create_bucket ? google_storage_bucket.main[0].name : var.bucket_name
  source              = data.archive_file.main.output_path
  content_disposition = "attachment"
  # TODO: The following blows up cloud build and gsutil cp downloads 0 bytes.
  # Even more surprising, this is not an issue with v1 functions - at least
  # cloud build works
  # content_encoding    = "gzip"
  content_type = "application/zip"
}

resource "google_cloudfunctions2_function" "main" {
  provider = google-beta
  name     = var.name
  project  = var.project_id

  location    = var.region
  description = var.description

  build_config {
    runtime               = var.runtime
    entry_point           = var.entry_point
    environment_variables = var.build_environment_variables
    source {
      storage_source {
        bucket = var.create_bucket ? google_storage_bucket.main[0].name : var.bucket_name
        object = google_storage_bucket_object.main.name
      }
    }
  }

  service_config {
    service_account_email          = var.service_account_email
    max_instance_count             = var.max_instances
    min_instance_count             = 1
    available_memory               = var.available_memory_mb # "256Mi"
    timeout_seconds                = var.timeout_s
    environment_variables          = var.environment_variables
    ingress_settings               = var.ingress_settings # "ALLOW_INTERNAL_ONLY"
    all_traffic_on_latest_revision = true
  }

  event_trigger {
    trigger               = var.trigger
    trigger_region        = var.region
    event_type            = var.event_type
    pubsub_topic          = var.pubsub_topic
    retry_policy          = var.retry_policy
    service_account_email = var.service_account_email
  }
}