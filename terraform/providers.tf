
terraform {
  required_version = "1.8.5"
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "5.32.0"
    }
  }
}

provider "google" {
  project = var.gcp_project_id
}
