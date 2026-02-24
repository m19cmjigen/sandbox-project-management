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
    key            = "production/vpc-iam.tfstate"
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
