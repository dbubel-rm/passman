variable "region" {
  description = "Region that the instances will be created"
  default = "us-east-1"
}

variable "production_database_name" {
  description = "The database name for staging"
  default = "passmandb"
}

variable "production_database_username" {
  description = "The username for the staging database"
  default = "passman"
}

variable "production_database_password" {
  description = "The user password for the staging database"
  default = "alkjswerfwere392kd"
}

variable "domain" {
  default = "The domain of your application"
}

variable "environment" {
  default = "dev"
}

variable "key_name" {
  description = "The public key for the bastion host"
  default = "production_key"
}
