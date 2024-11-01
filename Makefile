build:
	@go build -o bin/url_shortener

run: build
	@./bin/url_shortener

test:
	@go test ./... -v
