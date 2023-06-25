test:
	go test ./... -cover

gen-asset:
	go run tool/gen/gen_asset.go