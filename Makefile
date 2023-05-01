COVER_DIR=/tmp

test:
	go test -v -cover ./...

cover:
	go test -coverprofile=${COVER_DIR}/coverage.out ./...
	go tool cover -html=${COVER_DIR}/coverage.out
