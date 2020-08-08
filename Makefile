.PHONY: build
build:
	/usr/local/go/bin/go build -o detector -v ./cmd/detector/main.go && /usr/local/go/bin/go build -o daemon -v ./cmd/daemon/main.go

.DEFAULT_GOAL := build

.PHONY: build-daemon-run
build-daemon-run:
	/usr/local/go/bin/go build -o daemon -v ./cmd/daemon/main.go && sudo ./daemon

.PHONY: test
test:
	go test -v ./internal/app/store ./internal/app/tgpost ./internal/app/rpi-detector-mongo -cover

.PHONY: install
install:
	sudo ./install.sh
