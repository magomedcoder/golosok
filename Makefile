.PHONY: run
run:
	LD_LIBRARY_PATH=$(abspath ./build/golosok-linux/lib) \
	CGO_CPPFLAGS="-I$(abspath ./build/golosok-linux/lib)" \
	CGO_CFLAGS="-I$(abspath ./build/golosok-linux/lib)" \
	CGO_LDFLAGS="-L$(abspath ./build/golosok-linux/lib) -Wl,-rpath,'$$ORIGIN/../../build/golosok-linux/lib'" \
	go run -tags=cgo_vosk ./cmd/golosok -vosk-model ./build/golosok-linux/lib/models/vosk

.PHONY: build
build:
	CGO_ENABLED=1 go build -tags=cgo_vosk,portaudio -o build/golosok ./cmd/golosok

.PHONY: build-windows
build-windows:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o build/golosok.exe ./cmd/golosok
