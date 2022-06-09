locals {
  // Clumsy but works - and we could even put that bit in the root folder
  // only go.mod / go.sum and *.go
  excludes = [for s in fileset("${path.module}/fn", "*") : s if length(flatten(concat(regexall("go\\.(mod|sum)", s), regexall(".*go$", s)))) == 0]
}

output "path" {
  description = "The path to the function source"
  value       = "${path.module}/fn"
}

output "excludes" {
  description = "Files we want to exlude"
  value       = local.excludes
}

output "v1_entry_point" {
  description = "The v1 legacy label function entry point"
  value       = "LabelPubSub"
}

output "v2_entry_point" {
  description = "The v2 legacy label function entry point"
  value       = "LabelEvent"
}

output "entry_points_v1" {
  description = "The v1 function entry points provided by this module"
  value = {
    label           = "LabelPubSub"
    harden_sa       = "HardenPubSub"
    start_instances = "StartPubSub"
    stop_instances  = "StopPubSub"
  }
}

output "runtime" {
  description = "The runtime"
  value       = "go116"
}