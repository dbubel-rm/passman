variable "region" {
  description = "Region that the instances will be created"
}

variable "staging_database_name" {
  description = "The database name for staging"
}

variable "staging_database_username" {
  description = "The username for the staging database"
}

variable "staging_database_password" {
  description = "The user password for the staging database"
}

variable "domain" {
  default = "The domain of your application"
}

variable "environment" {
  default = "passman-staging"
}

variable "key_name" {
  description = "The public key for the bastion host"
}
