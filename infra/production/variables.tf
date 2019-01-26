variable "region" {
  description = "Region that the instances will be created"
}

variable "production_database_name" {
  description = "The database name for Production"
}

variable "production_database_username" {
  description = "The username for the Production database"
}

variable "production_database_password" {
  description = "The user password for the Production database"
}

variable "domain" {
  default = "The domain of your application"
}

variable "environment" {
  default = "passman-production"
}

variable "key_name" {
  description = "The public key for the bastion host"
}
