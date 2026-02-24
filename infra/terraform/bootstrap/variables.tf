variable "aws_region" {
  description = "AWS region where the state bucket and lock table will be created"
  type        = string
  default     = "ap-northeast-1"
}

variable "state_bucket_name" {
  description = "S3 bucket name for Terraform remote state"
  type        = string
  default     = "project-viz-terraform-state"
}

variable "lock_table_name" {
  description = "DynamoDB table name for Terraform state locking"
  type        = string
  default     = "project-viz-terraform-lock"
}
