run:
	go run main.go

run-jet:
	go run scene/jet/main.go

run-ui:
	go run scene/ui/main.go

test:
	go test ./... -cover

gen-asset:
	go run tool/gen_asset.go