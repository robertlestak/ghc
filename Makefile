build: clean dist-dirs build-darwin build-windows build-linux

VERSION=$$(git log --format="%h" -n 1 | tr -d '\n')

build-darwin: clean-darwin
	@ echo Building Version: $(VERSION)
	GOOS=darwin go build -ldflags "-X main.version=$(VERSION)" -o dist/ghc-darwin *.go

build-linux: clean-linux
	@ echo Building Version: $(VERSION)
	GOOS=linux go build -ldflags "-X main.version=$(VERSION)" -o dist/ghc-linux *.go

build-windows: clean-windows
	@ echo Building Version: $(VERSION)
	GOOS=windows go build -ldflags "-X main.version=$(VERSION)" -o dist/ghc-windows.exe *.go

dist-dirs:
	mkdir dist

clean-darwin:
	rm -rf dist/*-darwin

clean-linux:
	rm -rf dist/*-linux

clean-windows:
	rm -rf dist/*-windows*

clean:
	rm -rf dist

.PHONY: build build-darwin build-linux build-windows
.PHONY: dist-dirs clean
