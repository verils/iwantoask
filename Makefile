.PHONY: clean install build

clean:
	rm -f iwantoask*

install:
	go mod tidy

build:
	go build github.com/verils/iwantoask -o iwantoask
