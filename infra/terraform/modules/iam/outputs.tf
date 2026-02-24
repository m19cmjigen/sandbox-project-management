output "ecs_task_execution_role_arn" {
  description = "ARN of the ECS task execution role"
  value       = aws_iam_role.ecs_task_execution.arn
}

output "api_task_role_arn" {
  description = "ARN of the API ECS task role"
  value       = aws_iam_role.api_task.arn
}

output "batch_task_role_arn" {
  description = "ARN of the Batch ECS task role"
  value       = aws_iam_role.batch_task.arn
}

output "github_actions_role_arn" {
  description = "ARN of the GitHub Actions deploy role"
  value       = aws_iam_role.github_actions.arn
}
