
-- CREATE TABLE IF NOT EXISTS `users` (
--   `user_id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
--   `local_id` VARCHAR(32) NOT NULL,
--   `email` VARCHAR(32) UNIQUE NOT NULL,
--   `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
--   `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
--   `deleted_at` TIMESTAMP NULL,
--   PRIMARY KEY (`user_id`)
-- ) ENGINE=InnoDB DEFAULT CHARSET=utf8;
-- drop database if exists passman;
-- create database passman;
-- use passman;

CREATE TABLE IF NOT EXISTS `credentials` (
  `credential_id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `local_id` VARCHAR(29) NOT NULL,
  `service_name` VARCHAR(32) NOT NULL,
  `username` VARCHAR(128) NOT NULL,
  `password` VARCHAR(128) NOT NULL,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  -- `deleted_at` TIMESTAMP NULL,
  PRIMARY KEY (`credential_id`),
  UNIQUE KEY(`service_name`,`username`) 
  -- FOREIGN KEY (`user_id`) REFERENCES users(`user_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



