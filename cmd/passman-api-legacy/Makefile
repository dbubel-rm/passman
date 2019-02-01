
new_db:
	mysql -u root -e "drop database if exists passman"
	mysql -u root -e "create database passman"
	mysql -u root passman < schema.sql

deploy:
	$(AWS_ACCESS_KEY_ID=$EB_KEY AWS_SECRET_ACCESS_KEY=$EB_SECRET aws ecr get-login --no-include-email --region us-east-1)
	docker build -t passman-production-ecr .
	docker tag passman-production-ecr:latest 316188497159.dkr.ecr.us-east-1.amazonaws.com/passman-production-ecr:latest
	docker push 316188497159.dkr.ecr.us-east-1.amazonaws.com/passman-production-ecr:latest

.PHONY: new_db deploy