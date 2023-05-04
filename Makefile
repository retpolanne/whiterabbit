.PHONY: test

clean:
	@file pkg/whiterabbit.csv 2>/dev/null >/dev/null && rm pkg/whiterabbit.csv 2>/dev/null >/dev/null || true
	@ls out 2>/dev/null >/dev/null && rm -rf out

test: clean
	go test ./...

build: clean
	./build.sh
