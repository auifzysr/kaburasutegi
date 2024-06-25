terraform {
  backend "gcs" {
    bucket = "kaburasutegi-tfstate"
  }
}

variable "project" {
  type        = string
  description = "google cloud project"
}

variable "line_channel_token_secret_id" {
  type        = string
  description = "line_channel_token_secret_id"
}

variable "line_channel_secret_secret_id" {
  type        = string
  description = "line_channel_secret_secret_id"
}

locals {
  region                        = "asia-northeast1"
  project                       = var.project
  line_channel_token_secret_id  = var.line_channel_token_secret_id
  line_channel_secret_secret_id = var.line_channel_secret_secret_id
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

# TODO: upload sources to gcs first
# resource "google_cloudfunctions2_function" "callback" {
#   name     = "callback"
#   location = local.region
#   build_config {
#     runtime     = "go122"
#     entry_point = "entrypoint"
#   }
#   service_config {
#     available_memory   = 128
#     timeout_seconds    = 30
#     min_instance_count = 0
#     max_instance_count = 1
#     environment_variables = {
#       LOG_LEVEL = "debug"
#     }
#     service_account_email = google_service_account.function_runner.email
#     secret_environment_variables {
#       key        = "LINE_CHANNEL_TOKEN"
#       project_id = local.project
#       secret     = local.line_channel_token_secret_id
#       version    = "latest"
#     }
#     secret_environment_variables {
#       key        = "LINE_CHANNEL_SECRET"
#       project_id = local.project
#       secret     = local.line_channel_secret_secret_id
#       version    = "latest"
#     }
#     ingress_settings = "ALLOW_ALL"
#   }
# }
