terraform {
  required_version = ">= 1.6"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }

  # Remote state stored in S3 (bucket and table created by INFRA-002)
  backend "s3" {
    bucket         = "project-viz-terraform-state"
    key            = "production/main.tfstate"
    region         = "ap-northeast-1"
    dynamodb_table = "project-viz-terraform-lock"
    encrypt        = true
  }
}

provider "aws" {
  region = var.aws_region

  default_tags {
    tags = local.common_tags
  }
}

locals {
  project     = "project-viz"
  environment = "production"

  common_tags = {
    Project     = local.project
    Environment = local.environment
    ManagedBy   = "terraform"
  }
}

# ---------------------------------------------------------------------------
# VPC Module
# ---------------------------------------------------------------------------
module "vpc" {
  source = "../../modules/vpc"

  project     = local.project
  environment = local.environment
  vpc_cidr    = "10.1.0.0/16"

  public_subnets = {
    "a" = { cidr = "10.1.0.0/24", az = "${var.aws_region}a" }
    "c" = { cidr = "10.1.1.0/24", az = "${var.aws_region}c" }
  }

  private_subnets = {
    "a" = { cidr = "10.1.10.0/24", az = "${var.aws_region}a", nat_key = "a" }
    "c" = { cidr = "10.1.11.0/24", az = "${var.aws_region}c", nat_key = "c" }
  }

  # CloudWatch Logs: 90-day retention for production
  log_retention_days = 90
  tags               = local.common_tags
}

# ---------------------------------------------------------------------------
# ECS Module
# ---------------------------------------------------------------------------
module "ecs" {
  source = "../../modules/ecs"

  project     = local.project
  environment = local.environment
  aws_region  = var.aws_region

  vpc_id             = module.vpc.vpc_id
  public_subnet_ids  = module.vpc.public_subnet_ids
  private_subnet_ids = module.vpc.private_subnet_ids
  alb_sg_id          = module.vpc.alb_sg_id
  api_sg_id          = module.vpc.api_sg_id

  ecs_task_execution_role_arn = module.iam.ecs_task_execution_role_arn
  api_task_role_arn           = module.iam.api_task_role_arn
  batch_task_role_arn         = module.iam.batch_task_role_arn

  db_secret_arn        = module.aurora.db_secret_arn
  jwt_secret_arn       = var.jwt_secret_arn
  jira_secret_arn      = var.jira_secret_arn
  acm_certificate_arn  = var.acm_certificate_arn

  api_log_group_name   = module.vpc.api_log_group_name
  batch_log_group_name = module.vpc.batch_log_group_name

  # production: larger tasks, minimum 2 for zero-downtime deploy
  api_cpu           = 1024
  api_memory        = 2048
  api_desired_count = 2
  api_min_count     = 2
  api_max_count     = 8

  batch_cpu    = 1024
  batch_memory = 2048

  tags = local.common_tags
}

# ---------------------------------------------------------------------------
# CloudFront + S3 Module
# ---------------------------------------------------------------------------
module "cloudfront" {
  source = "../../modules/cloudfront"

  project     = local.project
  environment = local.environment

  # ACM certificate must be issued in us-east-1 for CloudFront
  acm_certificate_arn = var.cf_acm_certificate_arn
  domain_names        = var.cf_domain_names

  # PriceClass_All: global edge locations for production
  price_class = "PriceClass_All"

  tags = local.common_tags
}

# ---------------------------------------------------------------------------
# IAM Module
# ---------------------------------------------------------------------------
module "iam" {
  source = "../../modules/iam"

  project        = local.project
  environment    = local.environment
  aws_region     = var.aws_region
  aws_account_id = var.aws_account_id
  github_repo    = var.github_repo

  # OIDC provider is already created by staging; do not recreate.
  create_github_oidc_provider = false

  tags = local.common_tags
}

# ---------------------------------------------------------------------------
# Aurora Module
# ---------------------------------------------------------------------------
module "aurora" {
  source = "../../modules/aurora"

  project     = local.project
  environment = local.environment

  private_subnet_ids = module.vpc.private_subnet_ids
  db_sg_id           = module.vpc.db_sg_id

  # production: writer + reader, memory-optimized instance, 7-day backup, deletion protection
  instance_class        = "db.r8g.large"
  instance_count        = 2
  backup_retention_days = 7
  deletion_protection   = true

  tags = local.common_tags
}
