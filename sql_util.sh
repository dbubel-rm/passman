mysql -u root -e "drop database passman_test"
mysql -u root -e "create database passman_test"
mysql -u root passman_test < schema.sql