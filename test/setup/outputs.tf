/*
output "project_id" {
  value = module.project.project_id
}

output "sub_folder_id" {
  value = google_folder.ci_event_func_subfolder.id
}

output "sa_key" {
  value     = google_service_account_key.int_test.private_key
  sensitive = true
}

output "region" {
  value = var.region
}

output "zone" {
  value = var.zone
}

output "subnetwork" {
  value = module.network.subnets_self_links[0]
}
*/