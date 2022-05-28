output "project_id" {
  value       = var.project_id
  description = "The ID of the project to which resources are applied."
}

output "region" {
  value       = var.region
  description = "The region in which resources are applied."
}

output "function_name" {
  value       = module.pubsub.function_name
  description = "The name of the function."
}
