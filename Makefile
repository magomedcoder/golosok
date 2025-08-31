.PHONY: install
install:
	export LD_LIBRARY_PATH="$(pwd)/build/vosk:${LD_LIBRARY_PATH}"

.PHONY: build
build:
	CGO_ENABLED=1 go build -tags=cgo_vosk -o build/golosok ./cmd/golosok

.PHONY: run
run:
	go run -tags=cgo_vosk ./cmd/golosok
