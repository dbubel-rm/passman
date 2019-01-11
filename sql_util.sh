mysql -u root -e "drop database if exists passman"
mysql -u root -e "create database passman"


mysql -u root passman < schema.sql