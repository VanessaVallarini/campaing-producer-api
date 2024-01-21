provider "aws" {
    region = "us-east-1"
    access_key = "teste"
    secret_key = "teste"
    skip_credentials_validation = true
    skip_metadata_api_check     = true
    skip_requesting_account_id  = true
    endpoints {
        sqs = "http://localhost:4566"
  }
}

resource "aws_sqs_queue" "queue_campaing" {
  name = "queue_campaing"
}