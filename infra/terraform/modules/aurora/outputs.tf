output "cluster_endpoint" {
  description = "Writer endpoint of the Aurora cluster"
  value       = aws_rds_cluster.main.endpoint
}

output "cluster_reader_endpoint" {
  description = "Reader endpoint of the Aurora cluster"
  value       = aws_rds_cluster.main.reader_endpoint
}

output "cluster_identifier" {
  description = "Identifier of the Aurora cluster"
  value       = aws_rds_cluster.main.cluster_identifier
}

output "database_name" {
  description = "Name of the database"
  value       = aws_rds_cluster.main.database_name
}

output "db_secret_arn" {
  description = "ARN of the Secrets Manager secret containing DB credentials"
  value       = aws_secretsmanager_secret.db.arn
}

output "db_secret_name" {
  description = "Name of the Secrets Manager secret containing DB credentials"
  value       = aws_secretsmanager_secret.db.name
}
