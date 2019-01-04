mysql -u root -e "drop database passman"
mysql -u root -e "create database passman"


mysql -u root passman < schema.sql