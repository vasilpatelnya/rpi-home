.PHONY: build
build:
	/usr/local/go/bin/go build -o rpihome -v ./cmd/rpihome/rpihome.go

.DEFAULT_GOAL := build

.PHONY: build-run
build-run:
	/usr/local/go/bin/go build -o rpihome -v ./cmd/rpihome/rpihome.go && export ENVIRONMENT=default && ./rpihome

.PHONY: test
test:
	go test -p 1  ./... -v -coverprofile=coverage.out && go tool cover -func=coverage.out
	
.PHONY: install
install:
	sudo ./install.sh

gui:
	/usr/local/go/bin/go build -o gui ./cmd/gui && ./gui