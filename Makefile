.PHONY: build gen_example
build:
	go build -o build/ .

gen_example:
	cd example && buf generate .