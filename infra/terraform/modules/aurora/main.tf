terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.0"
    }
  }
}

# ---------------------------------------------------------------------------
# DB Subnet Group
# Aurora must be placed in a subnet group covering at least 2 AZs.
# ---------------------------------------------------------------------------
resource "aws_db_subnet_group" "main" {
  name        = "${var.project}-${var.environment}-aurora"
  description = "Subnet group for Aurora PostgreSQL cluster"
  subnet_ids  = var.private_subnet_ids

  tags = merge(var.tags, {
    Name = "${var.project}-${var.environment}-aurora-subnet-group"
  })
}

# ---------------------------------------------------------------------------
# DB Cluster Parameter Group
# ---------------------------------------------------------------------------
resource "aws_rds_cluster_parameter_group" "main" {
  name        = "${var.project}-${var.environment}-aurora-pg16"
  family      = "aurora-postgresql16"
  description = "Aurora PostgreSQL 16 cluster parameter group"

  # Force TLS connections
  parameter {
    name  = "rds.force_ssl"
    value = "1"
  }

  # Set timezone to JST
  parameter {
    name  = "timezone"
    value = "Asia/Tokyo"
  }

  tags = var.tags
}

# ---------------------------------------------------------------------------
# DB Instance Parameter Group
# ---------------------------------------------------------------------------
resource "aws_db_parameter_group" "main" {
  name        = "${var.project}-${var.environment}-aurora-pg16-instance"
  family      = "aurora-postgresql16"
  description = "Aurora PostgreSQL 16 instance parameter group"

  tags = var.tags
}

# ---------------------------------------------------------------------------
# Master Password (random, stored in Secrets Manager)
# ---------------------------------------------------------------------------
resource "random_password" "master" {
  length           = 32
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}

# ---------------------------------------------------------------------------
# Aurora Cluster
# ---------------------------------------------------------------------------
resource "aws_rds_cluster" "main" {
  cluster_identifier = "${var.project}-${var.environment}"
  engine             = "aurora-postgresql"
  engine_version     = var.engine_version

  database_name   = var.database_name
  master_username = var.master_username
  master_password = random_password.master.result

  db_subnet_group_name            = aws_db_subnet_group.main.name
  vpc_security_group_ids          = [var.db_sg_id]
  db_cluster_parameter_group_name = aws_rds_cluster_parameter_group.main.name

  # Encryption at rest
  storage_encrypted = true

  # Backup
  backup_retention_period   = var.backup_retention_days
  preferred_backup_window   = "17:00-18:00" # 02:00-03:00 JST
  preferred_maintenance_window = "sun:18:00-sun:19:00" # Sunday 03:00-04:00 JST

  # Deletion protection in production
  deletion_protection = var.deletion_protection

  # Skip final snapshot in non-production environments
  skip_final_snapshot       = !var.deletion_protection
  final_snapshot_identifier = var.deletion_protection ? "${var.project}-${var.environment}-final" : null

  # Enable enhanced monitoring logs
  enabled_cloudwatch_logs_exports = ["postgresql"]

  apply_immediately = !var.deletion_protection

  tags = merge(var.tags, {
    Name = "${var.project}-${var.environment}-aurora"
  })
}

# ---------------------------------------------------------------------------
# Aurora Instances
# ---------------------------------------------------------------------------
resource "aws_rds_cluster_instance" "main" {
  count = var.instance_count

  identifier         = "${var.project}-${var.environment}-${count.index}"
  cluster_identifier = aws_rds_cluster.main.id
  instance_class     = var.instance_class
  engine             = aws_rds_cluster.main.engine
  engine_version     = aws_rds_cluster.main.engine_version

  db_parameter_group_name = aws_db_parameter_group.main.name
  db_subnet_group_name    = aws_db_subnet_group.main.name

  # Enable enhanced monitoring (60-second intervals)
  monitoring_interval = 60
  monitoring_role_arn = aws_iam_role.rds_enhanced_monitoring.arn

  # Performance Insights
  performance_insights_enabled = true

  apply_immediately = !var.deletion_protection

  tags = merge(var.tags, {
    Name = "${var.project}-${var.environment}-aurora-${count.index}"
    Role = count.index == 0 ? "writer" : "reader"
  })
}

# ---------------------------------------------------------------------------
# Enhanced Monitoring IAM Role
# ---------------------------------------------------------------------------
resource "aws_iam_role" "rds_enhanced_monitoring" {
  name = "${var.project}-${var.environment}-rds-monitoring"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect    = "Allow"
      Principal = { Service = "monitoring.rds.amazonaws.com" }
      Action    = "sts:AssumeRole"
    }]
  })

  tags = var.tags
}

resource "aws_iam_role_policy_attachment" "rds_enhanced_monitoring" {
  role       = aws_iam_role.rds_enhanced_monitoring.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonRDSEnhancedMonitoringRole"
}

# ---------------------------------------------------------------------------
# Secrets Manager: DB credentials
# ECS tasks read this secret to get the connection string.
# ---------------------------------------------------------------------------
resource "aws_secretsmanager_secret" "db" {
  name        = "${var.project}/${var.environment}/db"
  description = "Aurora PostgreSQL credentials for ${var.project} ${var.environment}"

  # Automatically schedule deletion after 30 days (prevents accidental immediate deletion)
  recovery_window_in_days = 30

  tags = var.tags
}

resource "aws_secretsmanager_secret_version" "db" {
  secret_id = aws_secretsmanager_secret.db.id

  secret_string = jsonencode({
    username = var.master_username
    password = random_password.master.result
    host     = aws_rds_cluster.main.endpoint
    port     = 5432
    dbname   = var.database_name
    # Convenience: full DSN for golang-migrate and application use
    dsn = "postgres://${var.master_username}:${random_password.master.result}@${aws_rds_cluster.main.endpoint}:5432/${var.database_name}?sslmode=require"
  })
}

# ---------------------------------------------------------------------------
# CloudWatch Alarms
# ---------------------------------------------------------------------------
resource "aws_cloudwatch_metric_alarm" "db_cpu_high" {
  alarm_name          = "${var.project}-${var.environment}-aurora-cpu-high"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = 2
  metric_name         = "CPUUtilization"
  namespace           = "AWS/RDS"
  period              = 300
  statistic           = "Average"
  threshold           = 80
  alarm_description   = "Aurora CPU utilization exceeds 80%"

  dimensions = {
    DBClusterIdentifier = aws_rds_cluster.main.cluster_identifier
  }

  tags = var.tags
}

resource "aws_cloudwatch_metric_alarm" "db_freeable_memory_low" {
  alarm_name          = "${var.project}-${var.environment}-aurora-memory-low"
  comparison_operator = "LessThanThreshold"
  evaluation_periods  = 2
  metric_name         = "FreeableMemory"
  namespace           = "AWS/RDS"
  period              = 300
  statistic           = "Average"
  threshold           = 256 * 1024 * 1024 # 256 MB in bytes
  alarm_description   = "Aurora freeable memory is below 256 MB"

  dimensions = {
    DBClusterIdentifier = aws_rds_cluster.main.cluster_identifier
  }

  tags = var.tags
}
