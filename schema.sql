Create database passman;
CREATE TABLE IF NOT EXISTS `credentials` (
  `credential_id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `local_id` VARCHAR(29) NOT NULL,
  `service_name` VARCHAR(32) NOT NULL,
  `username` VARCHAR(128) NOT NULL,
  `password` VARCHAR(128) NOT NULL,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`credential_id`),
  UNIQUE KEY(`service_name`,`username`) 
  -- FOREIGN KEY (`user_id`) REFERENCES users(`user_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



-- CREATE USER ''@'ip-10-0-0-41' IDENTIFIED BY '';
-- GRANT SELECT on passman.credentials to 'passmanapp'@'%';
-- GRANT INSERT on passman.credentials to 'passmanapp'@'%';
-- GRANT DELETE on passman.credentials to 'passmanapp'@'%';
-- GRANT UPDATE on passman.credentials to 'passmanapp'@'%';