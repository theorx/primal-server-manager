test:
	go test ./... -race -coverprofile=cover.out && go tool cover -html=cover.out && rm cover.out