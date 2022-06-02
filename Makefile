run:
	go run main.go

test:
	go test ./...

gen-asset:
	go run tool/gen_asset.go