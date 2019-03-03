locals {
  staging_availability_zones = ["us-east-1a", "us-east-1b"]
}

provider "aws" {
  region = "${var.region}"
}

resource "aws_key_pair" "key" {
  key_name   = "${var.key_name}"
  public_key = "${file("keys/production_key.pub")}"
}

module "network_public_subnets" {
  source               = "../modules/network_public_subnets"
  environment          = "staging"
  vpc_cidr             = "10.0.0.0/16"
  public_subnets_cidr  = ["10.0.1.0/24", "10.0.2.0/24"]
#   private_subnets_cidr = ["10.0.10.0/24", "10.0.20.0/24"]
  region               = "${var.region}"
  availability_zones   = "${local.staging_availability_zones}"
  key_name             = "${var.key_name}"
}


module "rds" {
  source            = "../modules/rds"
  environment       = "dev"
  allocated_storage = "20"
  database_name     = "${var.production_database_name}"
  database_username = "${var.production_database_username}"
  database_password = "${var.production_database_password}"
  subnet_ids        = ["${module.network_public_subnets.public_subnet}"]
  vpc_id            = "${module.network_public_subnets.vpc_id}"
  instance_class    = "db.t2.micro"
}


module "ecs" {
  source      = "../modules/ecs"
  environment = "${var.environment}"
  vpc_id      = "${module.network_public_subnets.vpc_id}"

  # mysql_endpoint = "https://dynamodb.${var.region}.amazonaws.com"
  region = "${var.region}"

  # public_subnet_ids  = ["${module.network_public_subnets.public_subnets_id}"]
  public_subnet_ids = ["${module.network_public_subnets.public_subnet}"]

  # security_group_ids = ["${module.networking.default_sg_id}"]

  security_group_ids = [
    "${module.network_public_subnets.security_groups_ids}",
    "${module.rds.db_access_sg_id}",
  ]
  mysql_endpoint = "${module.rds.rds_address}"

  # database_name       = "${var.production_database_name}"
  # database_username   = "${var.production_database_username}"
  # database_password   = "${var.production_database_password}"
  # secret_key_base     = "${var.production_secret_key_base}"
}


module "bastion" {
  region      = "${var.region}"
  source      = "../modules/bastion"
  vpc_id      = "${module.network_public_subnets.vpc_id}"
  environment = "${var.environment}"

  security_group_ids = [
    "${module.network_public_subnets.bastion}",
    "${module.rds.db_access_sg_id}",
  ]

  public_subnet_ids = ["${module.network_public_subnets.public_subnet}"]
  key_name          = "production_key"
}

# module "ecs" {
#   source      = "../modules/ecs"
#   environment = "${var.environment}"
#   vpc_id      = "${module.networking.vpc_id}"

#   # mysql_endpoint = "https://dynamodb.${var.region}.amazonaws.com"
#   region = "${var.region}"

#   public_subnet_ids  = ["${module.networking.public_subnets_id}"]
#   private_subnet_ids = ["${module.networking.private_subnets_id}"]

#   # security_group_ids = ["${module.networking.default_sg_id}"]

#   security_group_ids = [
#     "${module.networking.security_groups_ids}",
#     "${module.rds.db_access_sg_id}",
#   ]
#   mysql_endpoint = "${module.rds.rds_address}"

#   # database_name       = "${var.staging_database_name}"
#   # database_username   = "${var.staging_database_username}"
#   # database_password   = "${var.staging_database_password}"
#   # secret_key_base     = "${var.staging_secret_key_base}"
# }
