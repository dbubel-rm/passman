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
    "us-east-1" = "ami-f652979b"
    "us-east-2" = "ami-fcc19b99"
    "us-west-1" = "ami-16efb076"
  }
}

variable "public_subnet_ids" {
  type        = "list"
  description = "asdf"
}
