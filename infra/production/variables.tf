variable "region" {
  description = "Region that the instances will be created"
}

/*====
environment specific variables
======*/

variable "production_database_name" {
  description = "The database name for Production"
}

variable "production_database_username" {
  description = "The username for the Production database"
}

variable "production_database_password" {
  description = "The user password for the Production database"
}

# variable "production_secret_key_base" {
#   description = "The Rails secret key for production"
# }

variable "domain" {
  default = "The domain of your application"
}

variable "environment" {
  default = "passman-production"
}


# Roles
# variable "role_arn" {
#   description = "The Role ARN URLs"
# }


variable "key_name" {
  description = "The public key for the bastion host"
}
