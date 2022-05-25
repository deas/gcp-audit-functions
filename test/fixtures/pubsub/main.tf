module "pubsub" {
  source = "../../../examples/pubsub"

  project_id = var.project_id
  # folder_id  = var.sub_folder_id
  organization_id = var.organization_id
  region          = var.region
}
