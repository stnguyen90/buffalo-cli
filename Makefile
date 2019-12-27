install:
	go install -v ./cmd/buffalo
	go mod tidy -v

test:
	go test -failfast -short -cover ./...
	go mod tidy -v

cov:
	go test -short -coverprofile cover.out ./...
	go tool cover -html cover.out
	go mod tidy -v
