variable "cluster_name" {
  description = "GKE Cluster Name"
  default     = "example"
}

variable "gcp_region" {
  description = "GCP region, e.g. us-east1"
  default     = "europe-west2"
}

variable "gcp_zone" {
  description = "GCP zone, e.g. us-east1-a"
  default     = "europe-west2-a"
}

variable "machine_type" {
  description = "GCP machine type"
  default     = "n1-standard-2"
}

