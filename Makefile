run:
	go mod tidy
	go run main.go

util-test:
	go test -v ./utils