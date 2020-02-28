latest:
	docker build -t passman .
	docker tag passman:latest stihl29/passman:latest
	docker push stihl29/passman:latest
build:
	GOOS=linux CGO_ENABLED=0 go build -v -o dist/passman .