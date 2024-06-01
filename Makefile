all:
	go build -o govhdx ./bin

generate:
	cd parser/ && binparsegen conversion.spec.yaml > vhdx_gen.go
