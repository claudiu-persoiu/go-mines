build:
	GOOS=js GOARCH=wasm go build -C ./wasm -o ../server/static/main.wasm

serve:
	go run ./server/server.go