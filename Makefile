build:
	@go build -o bin/golang-ex
run: build
	export DATABASE_URL="user=postgres dbname=postgres password=postgres sslmode=disable"
	@./bin/golang-ex
test: 
	@go test -v ./..