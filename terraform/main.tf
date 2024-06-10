data "google_service_account" "runtime" {
  project    = var.gcp_project_id
  account_id = "run-${var.app_name}"
}

resource "google_cloud_run_v2_service" "app" {
  name     = var.app_name
  project  = var.gcp_project_id
  location = var.region
  ingress  = "INGRESS_TRAFFIC_ALL"

  template {
    containers {
      image = var.image_tag

      ports {
        container_port = var.port
      }
      startup_probe {
        initial_delay_seconds = 5
        timeout_seconds       = 1
        period_seconds        = 3
        failure_threshold     = 1
        http_get {
          path = "/health"
        }
      }

      liveness_probe {
        http_get {
          path = "/health"
        }
      }

    }

    timeout = "5s"

    scaling {
      min_instance_count = 0
      max_instance_count = 1
    }

    service_account = data.google_service_account.runtime.email

  }
}

resource "google_cloud_run_v2_service_iam_member" "allow-public-acess" {
  location = google_cloud_run_v2_service.app.location
  name     = google_cloud_run_v2_service.app.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}

resource "google_cloud_run_domain_mapping" "apex-domain" {
  name     = var.domain
  location = google_cloud_run_v2_service.app.location
  metadata {
    namespace = var.gcp_project_id
  }
  spec {
    route_name = google_cloud_run_v2_service.app.name
  }
}

resource "google_cloud_run_domain_mapping" "www-domain" {
  name     = "www.${var.domain}"
  location = google_cloud_run_v2_service.app.location
  metadata {
    namespace = var.gcp_project_id
  }
  spec {
    route_name = google_cloud_run_v2_service.app.name
  }
}
