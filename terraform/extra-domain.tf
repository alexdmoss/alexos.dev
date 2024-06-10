resource "google_cloud_run_domain_mapping" "apex-alexmoss-dev" {
  name     = "alexmoss.dev"
  location = google_cloud_run_v2_service.app.location
  metadata {
    namespace = var.gcp_project_id
  }
  spec {
    route_name = google_cloud_run_v2_service.app.name
  }
}

resource "google_cloud_run_domain_mapping" "www-alexmoss-dev" {
  name     = "www.alexmoss.dev"
  location = google_cloud_run_v2_service.app.location
  metadata {
    namespace = var.gcp_project_id
  }
  spec {
    route_name = google_cloud_run_v2_service.app.name
  }
}
