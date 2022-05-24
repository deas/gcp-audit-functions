variable "project_id" {
  type        = string
  description = "The ID of the project to which resources will be applied."
}

variable "region" {
  type        = string
  description = "The region in which resources will be applied."
  default     = "us-west1" # Feature preview
}

variable "location" {
  type        = string
  description = "The region in which resources will be applied."
  default     = "US"
}