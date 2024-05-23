build-app:
	go build -o app -v
run-app:
	go build -o app -v ./...
	./app