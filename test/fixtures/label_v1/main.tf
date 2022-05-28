module "pubsub" {
  source = "../../../examples/label_v1"

  project_id = var.project_id
  # folder_id  = var.sub_folder_id
  org_id = var.org_id
  region = var.region
}
