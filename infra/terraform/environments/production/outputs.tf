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

output "alb_dns_name" {
  description = "ALB DNS name (configure Route 53 CNAME to this)"
  value       = module.ecs.alb_dns_name
}

output "api_ecr_repository_url" {
  description = "ECR repository URL for API image"
  value       = module.ecs.api_ecr_repository_url
}

output "batch_ecr_repository_url" {
  description = "ECR repository URL for batch image"
  value       = module.ecs.batch_ecr_repository_url
}

output "cloudfront_distribution_id" {
  description = "CloudFront distribution ID (for cache invalidation in CI/CD)"
  value       = module.cloudfront.distribution_id
}

output "cloudfront_domain_name" {
  description = "CloudFront domain name"
  value       = module.cloudfront.distribution_domain_name
}

output "frontend_s3_bucket" {
  description = "S3 bucket name for frontend assets"
  value       = module.cloudfront.s3_bucket_name
}
