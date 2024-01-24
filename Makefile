.PHONY: build clean deploy gomodgen

build: gomodgen
	export GO111MODULE=on
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/random random/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/song song/main.go

clean:
	rm -rf ./bin ./vendor go.sum

tidy: clean
	go mod tidy

deploy: clean tidy build
	sls deploy --verbose

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh
