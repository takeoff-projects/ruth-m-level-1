terraform {
  required_providers {
    google = {
      source = "hashicorp/google"
      version = "3.5.0"
    }
  }
}

provider "google" {
  credentials = file("~/roi-takeoff-key.json")

  project = "roi-takeoff-user47"
  region  = "us-central1"
  zone    = "us-central1-c"
}

resource "google_cloud_run_service" "events" {
  name     = "events-srv"
  location = "us-central1"

  template {
    spec {
      containers {
        image = "gcr.io/roi-takeoff-user47/events:v1.0"
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }
}

data "google_iam_policy" "noauth" {
  binding {
    role = "roles/run.invoker"
    members = [
      "allUsers",
    ]
  }
}

resource "google_cloud_run_service_iam_policy" "policy" {
  location = google_cloud_run_service.events.location
  project = google_cloud_run_service.events.project
  service = google_cloud_run_service.events.name
  policy_data = data.google_iam_policy.noauth.policy_data
}