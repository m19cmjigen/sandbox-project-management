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
    key            = "staging/vpc-iam.tfstate"
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
  environment = "staging"

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
  vpc_cidr    = "10.0.0.0/16"

  public_subnets = {
    "a" = { cidr = "10.0.0.0/24", az = "${var.aws_region}a" }
    "c" = { cidr = "10.0.1.0/24", az = "${var.aws_region}c" }
  }

  private_subnets = {
    "a" = { cidr = "10.0.10.0/24", az = "${var.aws_region}a", nat_key = "a" }
    "c" = { cidr = "10.0.11.0/24", az = "${var.aws_region}c", nat_key = "c" }
  }

  # CloudWatch Logs: 30-day retention for staging
  log_retention_days = 30
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

  # GitHub OIDC provider is shared across environments; create only once.
  create_github_oidc_provider = var.create_github_oidc_provider

  tags = local.common_tags
}
