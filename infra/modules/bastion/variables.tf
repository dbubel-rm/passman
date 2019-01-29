variable "vpc_id" {
  description = "The VPC ID"
}

variable "region" {
  description = "region"
}

variable "environment" {
  description = "The environment"
}

variable "key_name" {
  description = "name of key"
}

variable "security_group_ids" {
  type        = "list"
  description = "The security group collection"
}

variable "bastion_ami" {
  default = {
    "us-east-1" = "ami-0ac019f4fcb7cb7e6"
  }
}

variable "public_subnet_ids" {
  type        = "list"
  description = "public subnet ids"
}
