.PHONY: dev
dev:
	npx wrangler dev

.PHONY: build
build:
	go run github.com/syumai/workers/cmd/workers-assets-gen@latest -mode=go
	# tinygo build -panic=trap -gc=leaking -o ./build/app.wasm -target wasm -no-debug ./...
	GOOS=js GOARCH=wasm  go build -o ./build/app.wasm ./...

.PHONY: deploy
deploy:
	npx wrangler deploy
