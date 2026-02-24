variable "project" {
  description = "Project name used as a prefix for resource naming"
  type        = string
}

variable "environment" {
  description = "Deployment environment (staging or production)"
  type        = string
}

variable "acm_certificate_arn" {
  description = "ARN of the ACM certificate for CloudFront HTTPS. Must be in us-east-1 (CloudFront requirement)."
  type        = string
}

variable "domain_names" {
  description = "List of custom domain names (CNAMEs) for the CloudFront distribution"
  type        = list(string)
  default     = []
}

variable "price_class" {
  description = "CloudFront price class. PriceClass_100 covers US/EU/Asia."
  type        = string
  default     = "PriceClass_100"
  validation {
    condition     = contains(["PriceClass_100", "PriceClass_200", "PriceClass_All"], var.price_class)
    error_message = "price_class must be PriceClass_100, PriceClass_200, or PriceClass_All"
  }
}

variable "tags" {
  description = "Common tags applied to all resources"
  type        = map(string)
  default     = {}
}
