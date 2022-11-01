run:
	go run main.go

test:
	go test ./... -cover

gen-asset:
	go run tool/gen_asset.go