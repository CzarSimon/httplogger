test:
	go test ./...
	go vet ./...
	trivy fs .

image:
	sh build-image.sh

run:
	sh run-local.sh
