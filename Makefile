default:
	go build ./cmd/ufacter

test:
	go test ./...

update_modules:
	go get -d -u ./...
