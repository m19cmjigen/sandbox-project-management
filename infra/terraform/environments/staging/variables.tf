variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "ap-northeast-1"
}

variable "aws_account_id" {
  description = "AWS account ID"
  type        = string
}

variable "github_repo" {
  description = "GitHub repository in owner/repo format"
  type        = string
  default     = "m19cmjigen/sandbox-project-management"
}

variable "create_github_oidc_provider" {
  description = "Set to true when creating GitHub OIDC provider for the first time in this account"
  type        = bool
  default     = false
}

variable "acm_certificate_arn" {
  description = "ARN of the ACM certificate for HTTPS on ALB"
  type        = string
}

variable "jwt_secret_arn" {
  description = "ARN of the Secrets Manager secret for JWT_SECRET"
  type        = string
}

variable "jira_secret_arn" {
  description = "ARN of the Secrets Manager secret for Jira credentials"
  type        = string
}

variable "cf_acm_certificate_arn" {
  description = "ARN of the ACM certificate for CloudFront (must be in us-east-1)"
  type        = string
}

variable "cf_domain_names" {
  description = "Custom domain names for the CloudFront distribution"
  type        = list(string)
  default     = []
}
