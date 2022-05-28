variable "project_id" {
  type        = string
  description = "The ID of the project to which resources will be applied."
}

#variable "sub_folder_id" {
#  type        = string
#  description = "The ID of the folder to look for changes."
#}

variable "region" {
  type        = string
  description = "The region in which resources will be applied."
}

variable "org_id" {
  type        = string
  description = "The ID of the organization to which resources will be applied."
}
