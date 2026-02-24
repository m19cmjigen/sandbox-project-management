terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

# ---------------------------------------------------------------------------
# ECR Repositories
# ---------------------------------------------------------------------------
resource "aws_ecr_repository" "api" {
  name                 = "${var.project}-api"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }

  tags = merge(var.tags, {
    Name = "${var.project}-api"
  })
}

resource "aws_ecr_repository" "batch" {
  name                 = "${var.project}-batch"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }

  tags = merge(var.tags, {
    Name = "${var.project}-batch"
  })
}

# Lifecycle policy: keep only the latest 10 images to control storage cost
resource "aws_ecr_lifecycle_policy" "api" {
  repository = aws_ecr_repository.api.name

  policy = jsonencode({
    rules = [{
      rulePriority = 1
      description  = "Keep last 10 images"
      selection = {
        tagStatus   = "any"
        countType   = "imageCountMoreThan"
        countNumber = 10
      }
      action = { type = "expire" }
    }]
  })
}

resource "aws_ecr_lifecycle_policy" "batch" {
  repository = aws_ecr_repository.batch.name

  policy = jsonencode({
    rules = [{
      rulePriority = 1
      description  = "Keep last 10 images"
      selection = {
        tagStatus   = "any"
        countType   = "imageCountMoreThan"
        countNumber = 10
      }
      action = { type = "expire" }
    }]
  })
}

# ---------------------------------------------------------------------------
# ECS Cluster
# ---------------------------------------------------------------------------
resource "aws_ecs_cluster" "main" {
  name = "${var.project}-${var.environment}"

  setting {
    name  = "containerInsights"
    value = "enabled"
  }

  tags = merge(var.tags, {
    Name = "${var.project}-${var.environment}"
  })
}

resource "aws_ecs_cluster_capacity_providers" "main" {
  cluster_name       = aws_ecs_cluster.main.name
  capacity_providers = ["FARGATE", "FARGATE_SPOT"]

  default_capacity_provider_strategy {
    capacity_provider = "FARGATE"
    weight            = 1
  }
}

# ---------------------------------------------------------------------------
# Application Load Balancer
# ---------------------------------------------------------------------------
resource "aws_lb" "api" {
  name               = "${var.project}-${var.environment}-api"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [var.alb_sg_id]
  subnets            = var.public_subnet_ids

  # Enable access logs for security auditing
  access_logs {
    bucket  = "${var.project}-${var.environment}-alb-logs"
    enabled = false # Enable after creating the S3 bucket
  }

  tags = merge(var.tags, {
    Name = "${var.project}-${var.environment}-api-alb"
  })
}

# HTTP → HTTPS redirect
resource "aws_lb_listener" "http" {
  load_balancer_arn = aws_lb.api.arn
  port              = 80
  protocol          = "HTTP"

  default_action {
    type = "redirect"
    redirect {
      port        = "443"
      protocol    = "HTTPS"
      status_code = "HTTP_301"
    }
  }
}

# HTTPS listener (requires ACM certificate)
resource "aws_lb_listener" "https" {
  load_balancer_arn = aws_lb.api.arn
  port              = 443
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-TLS13-1-2-2021-06"
  certificate_arn   = var.acm_certificate_arn

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.api.arn
  }
}

# Target Group for API
resource "aws_lb_target_group" "api" {
  name        = "${var.project}-${var.environment}-api"
  port        = 8080
  protocol    = "HTTP"
  vpc_id      = var.vpc_id
  target_type = "ip" # required for Fargate

  health_check {
    enabled             = true
    path                = "/health"
    port                = "traffic-port"
    protocol            = "HTTP"
    healthy_threshold   = 2
    unhealthy_threshold = 3
    timeout             = 5
    interval            = 30
    matcher             = "200"
  }

  tags = merge(var.tags, {
    Name = "${var.project}-${var.environment}-api-tg"
  })
}

# ---------------------------------------------------------------------------
# ECS Task Definition: API
# ---------------------------------------------------------------------------
resource "aws_ecs_task_definition" "api" {
  family                   = "${var.project}-${var.environment}-api"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = var.api_cpu
  memory                   = var.api_memory
  execution_role_arn       = var.ecs_task_execution_role_arn
  task_role_arn            = var.api_task_role_arn

  container_definitions = jsonencode([{
    name      = "api"
    image     = "${aws_ecr_repository.api.repository_url}:latest"
    essential = true

    portMappings = [{
      containerPort = 8080
      protocol      = "tcp"
    }]

    environment = [
      { name = "PORT",      value = "8080" },
      { name = "GIN_MODE",  value = var.environment == "production" ? "release" : "debug" },
      { name = "LOG_FORMAT", value = "json" },
      { name = "LOG_LEVEL",  value = var.environment == "production" ? "info" : "debug" },
    ]

    # Inject credentials from Secrets Manager as environment variables
    secrets = [
      { name = "DB_HOST",     valueFrom = "${var.db_secret_arn}:host::" },
      { name = "DB_PORT",     valueFrom = "${var.db_secret_arn}:port::" },
      { name = "DB_USER",     valueFrom = "${var.db_secret_arn}:username::" },
      { name = "DB_PASSWORD", valueFrom = "${var.db_secret_arn}:password::" },
      { name = "DB_NAME",     valueFrom = "${var.db_secret_arn}:dbname::" },
      { name = "JWT_SECRET",  valueFrom = var.jwt_secret_arn },
    ]

    logConfiguration = {
      logDriver = "awslogs"
      options = {
        "awslogs-group"         = var.api_log_group_name
        "awslogs-region"        = var.aws_region
        "awslogs-stream-prefix" = "api"
      }
    }

    healthCheck = {
      command     = ["CMD-SHELL", "wget -qO- http://localhost:8080/health || exit 1"]
      interval    = 30
      timeout     = 5
      retries     = 3
      startPeriod = 60
    }
  }])

  tags = var.tags
}

# ---------------------------------------------------------------------------
# ECS Service: API
# ---------------------------------------------------------------------------
resource "aws_ecs_service" "api" {
  name            = "${var.project}-${var.environment}-api"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.api.arn
  desired_count   = var.api_desired_count
  launch_type     = "FARGATE"

  network_configuration {
    subnets          = var.private_subnet_ids
    security_groups  = [var.api_sg_id]
    assign_public_ip = false
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.api.arn
    container_name   = "api"
    container_port   = 8080
  }

  # Rolling update: keep minimum 50% healthy during deploy
  deployment_minimum_healthy_percent = 50
  deployment_maximum_percent         = 200

  deployment_circuit_breaker {
    enable   = true
    rollback = true
  }

  # Ignore task_definition changes so that CD pipeline can update the image
  # without Terraform detecting drift
  lifecycle {
    ignore_changes = [task_definition, desired_count]
  }

  depends_on = [aws_lb_listener.https]

  tags = var.tags
}

# ---------------------------------------------------------------------------
# Auto Scaling for API Service
# ---------------------------------------------------------------------------
resource "aws_appautoscaling_target" "api" {
  max_capacity       = var.api_max_count
  min_capacity       = var.api_min_count
  resource_id        = "service/${aws_ecs_cluster.main.name}/${aws_ecs_service.api.name}"
  scalable_dimension = "ecs:service:DesiredCount"
  service_namespace  = "ecs"
}

resource "aws_appautoscaling_policy" "api_cpu" {
  name               = "${var.project}-${var.environment}-api-cpu"
  policy_type        = "TargetTrackingScaling"
  resource_id        = aws_appautoscaling_target.api.resource_id
  scalable_dimension = aws_appautoscaling_target.api.scalable_dimension
  service_namespace  = aws_appautoscaling_target.api.service_namespace

  target_tracking_scaling_policy_configuration {
    predefined_metric_specification {
      predefined_metric_type = "ECSServiceAverageCPUUtilization"
    }
    target_value       = 70.0
    scale_in_cooldown  = 300
    scale_out_cooldown = 60
  }
}

# ---------------------------------------------------------------------------
# ECS Task Definition: Batch
# ---------------------------------------------------------------------------
resource "aws_ecs_task_definition" "batch" {
  family                   = "${var.project}-${var.environment}-batch"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = var.batch_cpu
  memory                   = var.batch_memory
  execution_role_arn       = var.ecs_task_execution_role_arn
  task_role_arn            = var.batch_task_role_arn

  container_definitions = jsonencode([{
    name      = "batch"
    image     = "${aws_ecr_repository.batch.repository_url}:latest"
    essential = true

    environment = [
      { name = "LOG_FORMAT",        value = "json" },
      { name = "LOG_LEVEL",         value = "info" },
      { name = "METRICS_NAMESPACE", value = "ProjectViz/${var.environment}" },
    ]

    secrets = [
      { name = "DB_HOST",       valueFrom = "${var.db_secret_arn}:host::" },
      { name = "DB_PORT",       valueFrom = "${var.db_secret_arn}:port::" },
      { name = "DB_USER",       valueFrom = "${var.db_secret_arn}:username::" },
      { name = "DB_PASSWORD",   valueFrom = "${var.db_secret_arn}:password::" },
      { name = "DB_NAME",       valueFrom = "${var.db_secret_arn}:dbname::" },
      { name = "JIRA_BASE_URL", valueFrom = "${var.jira_secret_arn}:base_url::" },
      { name = "JIRA_EMAIL",    valueFrom = "${var.jira_secret_arn}:email::" },
      { name = "JIRA_API_TOKEN", valueFrom = "${var.jira_secret_arn}:api_token::" },
    ]

    logConfiguration = {
      logDriver = "awslogs"
      options = {
        "awslogs-group"         = var.batch_log_group_name
        "awslogs-region"        = var.aws_region
        "awslogs-stream-prefix" = "batch"
      }
    }
  }])

  tags = var.tags
}

# ---------------------------------------------------------------------------
# EventBridge Scheduled Rules for Batch
# ---------------------------------------------------------------------------

# Full Sync: daily at 02:00 JST (17:00 UTC)
resource "aws_cloudwatch_event_rule" "batch_full_sync" {
  name                = "${var.project}-${var.environment}-batch-full-sync"
  description         = "Trigger full sync batch daily at 02:00 JST"
  schedule_expression = "cron(0 17 * * ? *)"

  tags = var.tags
}

resource "aws_cloudwatch_event_target" "batch_full_sync" {
  rule      = aws_cloudwatch_event_rule.batch_full_sync.name
  target_id = "BatchFullSync"
  arn       = aws_ecs_cluster.main.arn
  role_arn  = aws_iam_role.eventbridge_ecs.arn

  ecs_target {
    task_count          = 1
    task_definition_arn = aws_ecs_task_definition.batch.arn
    launch_type         = "FARGATE"
    platform_version    = "LATEST"

    network_configuration {
      subnets          = var.private_subnet_ids
      security_groups  = [var.api_sg_id]
      assign_public_ip = false
    }
  }

  input = jsonencode({
    containerOverrides = [{
      name    = "batch"
      environment = [
        { name = "BATCH_SYNC_MODE", value = "full" }
      ]
    }]
  })
}

# Delta Sync: every hour at :00
resource "aws_cloudwatch_event_rule" "batch_delta_sync" {
  name                = "${var.project}-${var.environment}-batch-delta-sync"
  description         = "Trigger delta sync batch every hour"
  schedule_expression = "cron(0 * * * ? *)"

  tags = var.tags
}

resource "aws_cloudwatch_event_target" "batch_delta_sync" {
  rule      = aws_cloudwatch_event_rule.batch_delta_sync.name
  target_id = "BatchDeltaSync"
  arn       = aws_ecs_cluster.main.arn
  role_arn  = aws_iam_role.eventbridge_ecs.arn

  ecs_target {
    task_count          = 1
    task_definition_arn = aws_ecs_task_definition.batch.arn
    launch_type         = "FARGATE"
    platform_version    = "LATEST"

    network_configuration {
      subnets          = var.private_subnet_ids
      security_groups  = [var.api_sg_id]
      assign_public_ip = false
    }
  }

  input = jsonencode({
    containerOverrides = [{
      name    = "batch"
      environment = [
        { name = "BATCH_SYNC_MODE", value = "delta" }
      ]
    }]
  })
}

# IAM role for EventBridge to run ECS tasks
resource "aws_iam_role" "eventbridge_ecs" {
  name = "${var.project}-${var.environment}-eventbridge-ecs"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect    = "Allow"
      Principal = { Service = "events.amazonaws.com" }
      Action    = "sts:AssumeRole"
    }]
  })

  tags = var.tags
}

resource "aws_iam_role_policy" "eventbridge_ecs" {
  name = "RunECSTask"
  role = aws_iam_role.eventbridge_ecs.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect   = "Allow"
        Action   = ["ecs:RunTask"]
        Resource = [aws_ecs_task_definition.batch.arn]
      },
      {
        Effect   = "Allow"
        Action   = ["iam:PassRole"]
        Resource = [var.ecs_task_execution_role_arn, var.batch_task_role_arn]
      },
    ]
  })
}
