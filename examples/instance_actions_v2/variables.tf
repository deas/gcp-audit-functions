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

variable "time_zone" {
  type        = string
  description = "The timezone to use in scheduler"
  default     = "Etc/UTC"
}

variable "search_scope" {
  type        = string
  description = "The scope of the search"
  default     = "projects"
}

variable "action" {
  type        = map(any)
  description = "Instance action parameters"
  default = {
    "start" = {
      "schedule" = "0 1 * * *"
      "query"    = "labels.start_daily:true AND state:TERMINATED"
    }
    "stop" = {
      "schedule" = "0 2 * * *"
      "query" : "labels.stop_daily:true AND state:RUNNING"
    }
  }
}
variable "service_account_email" {
  type        = string
  description = "The service account email"
  default     = ""
}