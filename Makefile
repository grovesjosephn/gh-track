.PHONY: build run clean

build:
	go build -o hab main.go

run:
	go run main.go

clean:
	rm -f hab

install: build
	cp hab /usr/local/bin/