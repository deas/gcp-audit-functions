resource "random_pet" "main" {
  separator = "-"
}


resource "google_compute_instance" "main" {
  boot_disk {
    initialize_params {
      image = var.image
    }
  }

  machine_type = var.machine_type
  name         = "unlabelled-${random_pet.main.id}"
  zone         = var.zone

  lifecycle {
    ignore_changes = [
      labels,
    ]
  }

  network_interface {
    subnetwork = var.subnetwork
  }

  project = var.project_id
}