.PHONY: start, test, doc

start:
	go run cmd/app/main.go 
test: 
	go test -v ./... -run=.

doc: 
	cd cmd/app; swag init --parseDependency --parseInternal --parseDepth 1 
