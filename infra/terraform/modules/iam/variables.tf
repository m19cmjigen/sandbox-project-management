variable "project" {
  description = "Project name used as a prefix for resource naming"
  type        = string
}

variable "environment" {
  description = "Deployment environment (staging or production)"
  type        = string
}

variable "aws_region" {
  description = "AWS region where resources are deployed"
  type        = string
}

variable "aws_account_id" {
  description = "AWS account ID (used for ARN construction)"
  type        = string
}

variable "github_repo" {
  description = "GitHub repository in owner/repo format (used for OIDC trust policy)"
  type        = string
}

variable "create_github_oidc_provider" {
  description = "Set to true to create the GitHub OIDC provider. Only one provider can exist per account; set to false if already created."
  type        = bool
  default     = false
}

variable "tags" {
  description = "Common tags applied to all resources"
  type        = map(string)
  default     = {}
}
