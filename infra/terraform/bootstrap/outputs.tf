output "state_bucket_name" {
  description = "S3 bucket name to use in environment backend configurations"
  value       = aws_s3_bucket.state.bucket
}

output "state_bucket_arn" {
  description = "ARN of the Terraform state S3 bucket"
  value       = aws_s3_bucket.state.arn
}

output "lock_table_name" {
  description = "DynamoDB table name to use in environment backend configurations"
  value       = aws_dynamodb_table.lock.name
}
