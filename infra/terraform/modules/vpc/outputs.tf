output "vpc_id" {
  description = "ID of the VPC"
  value       = aws_vpc.main.id
}

output "vpc_cidr" {
  description = "CIDR block of the VPC"
  value       = aws_vpc.main.cidr_block
}

output "public_subnet_ids" {
  description = "List of public subnet IDs"
  value       = [for s in aws_subnet.public : s.id]
}

output "private_subnet_ids" {
  description = "List of private subnet IDs"
  value       = [for s in aws_subnet.private : s.id]
}

output "alb_sg_id" {
  description = "Security group ID for the Application Load Balancer"
  value       = aws_security_group.alb.id
}

output "api_sg_id" {
  description = "Security group ID for the Backend API ECS tasks"
  value       = aws_security_group.api.id
}

output "db_sg_id" {
  description = "Security group ID for Aurora PostgreSQL"
  value       = aws_security_group.db.id
}

output "api_log_group_name" {
  description = "CloudWatch log group name for API service"
  value       = aws_cloudwatch_log_group.api.name
}

output "batch_log_group_name" {
  description = "CloudWatch log group name for Batch service"
  value       = aws_cloudwatch_log_group.batch.name
}
