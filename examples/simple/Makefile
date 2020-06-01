build:
	mkdir -p bin
	go build -o bin/protoc-gen-simple
.PHONY: build

run: build
	rm -f *.py
	protoc -I .  --plugin=./bin/protoc-gen-simple --python_out=. --simple_out=. *.proto
	python3 simple.decode.pb.py
.PHONY: run