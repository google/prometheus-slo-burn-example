terraform {
  required_version = ">= 0.11.1"
}

provider "google" {
  region = var.gcp_region
}

data "google_client_config" "current" {
}

resource "google_container_cluster" "example-cluster" {
  name               = var.cluster_name
  description        = "prometheus example k8s cluster"
  region             = var.gcp_region
  initial_node_count = "1"

  logging_service    = "logging.googleapis.com/kubernetes"
  monitoring_service = "monitoring.googleapis.com/kubernetes"

  // Use legacy ABAC until these issues are resolved:
  //   https://github.com/mcuadros/terraform-provider-helm/issues/56
  //   https://github.com/terraform-providers/terraform-provider-kubernetes/pull/73
  enable_legacy_abac = true

  remove_default_node_pool = true
}

resource "google_container_node_pool" "pool0" {
  name       = "pool-0"
  cluster    = google_container_cluster.example-cluster.name
  node_count = 1
  region     = var.gcp_region

  autoscaling {
    min_node_count = 1
    max_node_count = 5
  }

  management {
    auto_repair  = "true"
    auto_upgrade = "true"
  }

  node_config {
    machine_type = var.machine_type
    preemptible  = "true"

    metadata = {
      disable-legacy-endpoints = "true"
    }

    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform",
      "https://www.googleapis.com/auth/cloud_debugger",
      "https://www.googleapis.com/auth/compute",
      "https://www.googleapis.com/auth/devstorage.read_only",
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring",
      "https://www.googleapis.com/auth/service.management",
      "https://www.googleapis.com/auth/servicecontrol",
      "https://www.googleapis.com/auth/source.read_only",
      "https://www.googleapis.com/auth/taskqueue",
      "https://www.googleapis.com/auth/trace.append",
    ]
  }
}

resource "google_compute_global_address" "grafana-ip" {
  name = "grafana-ip"
}

resource "google_compute_global_address" "server-ip" {
  name = "server-ip"
}

