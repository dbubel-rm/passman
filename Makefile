build:
	go build -ldflags "\
              -X main.BUILD_GIT_HASH=`git rev-parse HEAD` \
              -X main.BUILD_DATE=`date -u +'%Y-%m-%dT%H:%M:%SZ'`" \
              -v -o passman main.go

test:
	docker-compose -f docker-compose.yaml up --build --abort-on-container-exit
test-local:
	mysql -u root -e "create database if not exists passman;"
	cd cmd/passman-api/ && MYSQL_ENDPOINT=localhost go test -v ./...
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
	docker-compose -f docker-compose.prod.yaml up --scale passman-api=2--build -d
stop:
	docker-compose -f docker-compose.prod.yaml down

.PHONY: test test-local run build deploy run-dev stop-dev start stop
