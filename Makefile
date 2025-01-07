test:
	go test -timeout 5s -race ./...
coverage:
	go test -timeout 5s -race ./... -coverprofile=cover.out
	go tool cover -html=cover.out