terraform {
  backend "gcs" {
    bucket = "kaburasutegi-tfstate"
  }
}

variable "project" {
  type        = string
  description = "google cloud project"
}

locals {
  region  = "asia-northeast1"
  project = var.project
}

provider "google" {
  project = local.project
  region  = local.region
}

resource "google_service_account" "function_runner" {
  project      = local.project
  account_id   = "function-runner"
  display_name = "function-runner"
}

resource "google_project_iam_member" "secretmanager_access" {
  project = local.project
  role    = "roles/secretmanager.secretAccessor"
  member  = google_service_account.function_runner.member

  depends_on = [google_service_account.function_runner]
}
