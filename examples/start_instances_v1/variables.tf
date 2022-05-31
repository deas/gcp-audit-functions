variable "project_id" {
  type        = string
  description = "The ID of the project to which resources will be applied."
}

variable "region" {
  type        = string
  description = "The region in which resources will be applied."
}

variable "org_id" {
  type        = string
  description = "The organization ID to which resources will be applied."
  default     = "override in terraform.tfvars"
}

variable "vm" {
  type = object({
    zone       = string
    subnetwork = string
  })
  description = "VM spec - zone and subnetwork. Null to disable"
  default     = null
}

variable "search" {
  type        = string
  description = "The asset search"
  default     = <<EOF
{
  "scope": "organizations/your-org-id",
  "query": "labels.start_daily:true AND state:TERMINATED",
  "assetTypes": ["compute.googleapis.com/Instance"]
}
EOF
}

variable "schedule" {
  type        = string
  description = "The schedule"
  default     = "0 1 * * *"
}

variable "service_account_email" {
  type        = string
  description = "The service account email"
  default     = ""
}