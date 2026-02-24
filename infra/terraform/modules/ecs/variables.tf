variable "project" {
  description = "Project name used as a prefix for resource naming"
  type        = string
}

variable "environment" {
  description = "Deployment environment (staging or production)"
  type        = string
}

variable "aws_region" {
  description = "AWS region"
  type        = string
}

variable "vpc_id" {
  description = "VPC ID (from vpc module)"
  type        = string
}

variable "public_subnet_ids" {
  description = "Public subnet IDs for ALB"
  type        = list(string)
}

variable "private_subnet_ids" {
  description = "Private subnet IDs for ECS tasks"
  type        = list(string)
}

variable "alb_sg_id" {
  description = "Security group ID for ALB (from vpc module)"
  type        = string
}

variable "api_sg_id" {
  description = "Security group ID for API ECS tasks (from vpc module)"
  type        = string
}

variable "ecs_task_execution_role_arn" {
  description = "ARN of the ECS task execution role (from iam module)"
  type        = string
}

variable "api_task_role_arn" {
  description = "ARN of the API ECS task role (from iam module)"
  type        = string
}

variable "batch_task_role_arn" {
  description = "ARN of the batch ECS task role (from iam module)"
  type        = string
}

variable "db_secret_arn" {
  description = "ARN of the Secrets Manager secret containing DB credentials (from aurora module)"
  type        = string
}

variable "jwt_secret_arn" {
  description = "ARN of the Secrets Manager secret containing JWT_SECRET"
  type        = string
}

variable "jira_secret_arn" {
  description = "ARN of the Secrets Manager secret containing Jira credentials"
  type        = string
}

variable "api_log_group_name" {
  description = "CloudWatch log group name for API (from vpc module)"
  type        = string
}

variable "batch_log_group_name" {
  description = "CloudWatch log group name for batch (from vpc module)"
  type        = string
}

variable "acm_certificate_arn" {
  description = "ARN of the ACM certificate for HTTPS on ALB"
  type        = string
}

# API task sizing
variable "api_cpu" {
  description = "CPU units for API task (256 = 0.25 vCPU)"
  type        = number
  default     = 512
}

variable "api_memory" {
  description = "Memory (MB) for API task"
  type        = number
  default     = 1024
}

variable "api_desired_count" {
  description = "Initial desired number of API tasks"
  type        = number
  default     = 1
}

variable "api_min_count" {
  description = "Minimum number of API tasks for Auto Scaling"
  type        = number
  default     = 1
}

variable "api_max_count" {
  description = "Maximum number of API tasks for Auto Scaling"
  type        = number
  default     = 4
}

# Batch task sizing
variable "batch_cpu" {
  description = "CPU units for batch task"
  type        = number
  default     = 512
}

variable "batch_memory" {
  description = "Memory (MB) for batch task"
  type        = number
  default     = 1024
}

variable "tags" {
  description = "Common tags applied to all resources"
  type        = map(string)
  default     = {}
}
