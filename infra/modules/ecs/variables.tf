variable "vpc_id" {
  description = "The VPC ID"
}

variable "environment" {
  description = "The environment"
}

variable "region" {
  description = "The AWS region"
}

variable "mysql_endpoint" {
  description = "The mysql endpoint"
}

# variable "role_arn" {
#     description = "The ARN for roles"
# }

variable "public_subnet_ids" {
  type        = "list"
  description = "The public subnet IDs"
}

# variable "private_subnet_ids" {
#   type        = "list"
#   description = "The private subnet IDs"
# }

variable "security_group_ids" {
  type        = "list"
  description = "The security group collection"
}
