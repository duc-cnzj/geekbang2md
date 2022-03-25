.PHONY: lint
lint:
	golangci-lint run -D errcheck

.PHONY: fmt
fmt:
	gofmt -s -w ./ && goimports -w ./

.PHONY: build_race
build_race:
	CGO_ENABLED=1 go build -ldflags="-w -s" -race -o geekbang2md .

.PHONY: build
build:
	CGO_ENABLED=0 go build -ldflags="-w -s" -o geekbang2md .

.PHONY: build_drawin_amd64
build_drawin_amd64:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-w -s" -o geekbang2md-darwin-amd64 .

.PHONY: build_drawin_arm64
build_drawin_arm64:
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="-w -s" -o geekbang2md-darwin-arm64 .

.PHONY: build_windows
build_windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags='-w -s' -o geekbang2md-windows.exe .