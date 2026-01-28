build:
	GOOS=js GOARCH=wasm go build -C ./wasm -o ../server/static/main.wasm

build-tiny:
	cd wasm; tinygo build -target wasm -o ../server/static/main.wasm -no-debug; cd -

serve:
	go run -C server/ server.go
