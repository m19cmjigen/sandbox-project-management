output "cluster_name" {
  description = "ECS cluster name"
  value       = aws_ecs_cluster.main.name
}

output "cluster_arn" {
  description = "ECS cluster ARN"
  value       = aws_ecs_cluster.main.arn
}

output "api_service_name" {
  description = "ECS service name for the API"
  value       = aws_ecs_service.api.name
}

output "api_ecr_repository_url" {
  description = "ECR repository URL for the API image"
  value       = aws_ecr_repository.api.repository_url
}

output "batch_ecr_repository_url" {
  description = "ECR repository URL for the batch image"
  value       = aws_ecr_repository.batch.repository_url
}

output "alb_dns_name" {
  description = "DNS name of the Application Load Balancer"
  value       = aws_lb.api.dns_name
}

output "alb_zone_id" {
  description = "Hosted zone ID of the ALB (for Route 53 alias records)"
  value       = aws_lb.api.zone_id
}

output "api_task_definition_arn" {
  description = "ARN of the API task definition"
  value       = aws_ecs_task_definition.api.arn
}

output "batch_task_definition_arn" {
  description = "ARN of the batch task definition"
  value       = aws_ecs_task_definition.batch.arn
}
