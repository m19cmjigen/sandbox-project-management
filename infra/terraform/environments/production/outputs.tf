output "vpc_id" {
  description = "VPC ID"
  value       = module.vpc.vpc_id
}

output "public_subnet_ids" {
  description = "Public subnet IDs"
  value       = module.vpc.public_subnet_ids
}

output "private_subnet_ids" {
  description = "Private subnet IDs"
  value       = module.vpc.private_subnet_ids
}

output "alb_sg_id" {
  description = "ALB security group ID"
  value       = module.vpc.alb_sg_id
}

output "api_sg_id" {
  description = "API security group ID"
  value       = module.vpc.api_sg_id
}

output "db_sg_id" {
  description = "DB security group ID"
  value       = module.vpc.db_sg_id
}

output "ecs_task_execution_role_arn" {
  description = "ECS task execution role ARN"
  value       = module.iam.ecs_task_execution_role_arn
}

output "api_task_role_arn" {
  description = "API task role ARN"
  value       = module.iam.api_task_role_arn
}

output "batch_task_role_arn" {
  description = "Batch task role ARN"
  value       = module.iam.batch_task_role_arn
}

output "github_actions_role_arn" {
  description = "GitHub Actions deploy role ARN"
  value       = module.iam.github_actions_role_arn
}

output "aurora_cluster_endpoint" {
  description = "Aurora cluster writer endpoint"
  value       = module.aurora.cluster_endpoint
}

output "aurora_db_secret_arn" {
  description = "Secrets Manager ARN for DB credentials"
  value       = module.aurora.db_secret_arn
}
