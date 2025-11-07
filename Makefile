build_linux:
	go env -w CGO_ENABLED=0 GOOS=linux && \
	go build -o invoice-api ./cmd/api
build_windows:
	go env -w  CGO_ENABLED=0 GOOS=windows GOARCH=amd64 && \
	go build -o invoice-api.exe ./cmd/api
build_mac_amd64:
	go env -w CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 &&\
	go build -o invoice-api-mac-amd64 ./cmd/api
build_mac_arm64:
	go env -w CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 && \
	go build -o invoice-api-mac-arm64 ./cmd/api

run_dev:
	make build_linux && docker-compose build && docker-compose up