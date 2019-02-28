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
	cd cmd/passman-api/ && DB_HOST="root:my-secret-pw@tcp(localhost:3306)/passman" go test -v ./...
run:
	docker-compose -f docker-compose.prod.yaml up --build -d
stop:
	docker-compose -f docker-compose.prod.yaml down

.PHONY: build
