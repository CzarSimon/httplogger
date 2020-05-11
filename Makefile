test:
	go test ./...

image:
	sh build-image.sh

run:
	sh run-local.sh
