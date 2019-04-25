build:
	docker build \
		-t sales-api-amd64:1.0 \
		--build-arg PACKAGE_NAME=sales-api \
		--build-arg VCS_REF=`git rev-parse HEAD` \
		--build-arg BUILD_DATE=`date -u +”%Y-%m-%dT%H:%M:%SZ”` \
		.
	docker system prune -f
	
test:
	docker-compose -f docker-compose.yaml up --build --abort-on-container-exit 
test-local:
	mysql -u root -e "create database if not exists passman;"
	cd cmd/passman-api/ && DB_HOST="root@tcp(127.0.0.1:3306)/passman" go test -v ./...
run-dev:
	docker-compose -f docker-compose.dev.yaml up --build -d
stop-dev:
	docker-compose -f docker-compose.dev.yaml down
deploy:
	# $(AWS_ACCESS_KEY_ID=$EB_KEY AWS_SECRET_ACCESS_KEY=$EB_SECRET aws ecr get-login --no-include-email --region us-east-1)
	docker build -t passman .
	docker tag passman:latest stihl29/passman:latest
	docker push stihl29/passman:latest
cli:
	cd cmd/passman-cli && CGO_ENABLED=0 go build -i -a -v -o `$GOPATH/bin/passman` -ldflags '-extldflags "-static"' main.go
start:
	MYSQL_ENDPOINT=${MYSQL_ENDPOINT} \
	MYSQL_USERNAME=${MYSQL_USERNAME} \
	MYSQL_PASSWORD=${MYSQL_PASSWORD} \
	MYSQL_DB=${MYSQL_DB} \
	docker-compose -f docker-compose.prod.yaml up --build -d
stop: 
	docker-compose -f docker-compose.prod.yaml down

.PHONY: test test-local run build deploy run-dev stop-dev start stop
