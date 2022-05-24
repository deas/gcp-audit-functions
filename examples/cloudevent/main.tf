resource "random_pet" "main" {
  separator = "-"
}

# TODO: A bit hacky? Open for ideas.
module "function" {
  source = "../.."
}

module "audit_label" {
  source = "../../modules/functions-v2"

  description = "Labels resource with owner information."
  entry_point = module.function.v2_entry_point

  environment_variables = {
    LABEL_KEY = "principal-email"
  }

  # event_trigger                  = module.event_log_entry.function_event_trigger
  name                           = "audit-label-${random_pet.main.id}"
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