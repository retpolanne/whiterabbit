.PHONY: test

clean:
	@file pkg/whiterabbit.csv 2>/dev/null >/dev/null && rm pkg/whiterabbit.csv 2>/dev/null >/dev/null || true

test: clean
	go test ./...

