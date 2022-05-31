variable "project_id" {
  type        = string
  description = "The ID of the project to which resources will be applied."
}

variable "zone" {
  type        = string
  description = "The zone in which resources will be applied."
}

variable "subnetwork" {
  type        = string
  description = "The name or self_link of the subnetwork to create compute instance in."
}

variable "image" {
  type        = string
  default     = "debian-cloud/debian-9"
  description = "The image to use for the compute instance."
}

# https://cloud.google.com/compute/docs/machine-types
variable "machine_type" {
  type        = string
  default     = "f1-micro"
  description = "The machine type to use for the compute instance."

}