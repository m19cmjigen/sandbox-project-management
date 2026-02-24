variable "project" {
  description = "Project name used as a prefix for resource naming"
  type        = string
}

variable "environment" {
  description = "Deployment environment (staging or production)"
  type        = string
}

variable "private_subnet_ids" {
  description = "List of private subnet IDs where Aurora will be placed"
  type        = list(string)
}

variable "db_sg_id" {
  description = "Security group ID for Aurora (from VPC module)"
  type        = string
}

variable "engine_version" {
  description = "Aurora PostgreSQL engine version"
  type        = string
  default     = "16.4"
}

variable "database_name" {
  description = "Name of the initial database to create"
  type        = string
  default     = "project_visualization"
}

variable "master_username" {
  description = "Master DB username"
  type        = string
  default     = "dbadmin"
}

variable "instance_class" {
  description = "Instance class for Aurora cluster instances"
  type        = string
}

variable "instance_count" {
  description = "Number of Aurora instances (1 = writer only, 2+ = writer + readers)"
  type        = number
  default     = 1
  validation {
    condition     = var.instance_count >= 1
    error_message = "instance_count must be at least 1"
  }
}

variable "backup_retention_days" {
  description = "Number of days to retain automated backups"
  type        = number
  default     = 7
}

variable "deletion_protection" {
  description = "Enable deletion protection (set true for production)"
  type        = bool
  default     = false
}

variable "tags" {
  description = "Common tags applied to all resources"
  type        = map(string)
  default     = {}
}
