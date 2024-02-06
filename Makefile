.PHONY: build clean gosumgen

build: gosumgen
	export GO111MODULE=on
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/random random/main.go

clean:
	rm -rf ./bin ./vendor go.sum
	
gosumgen: clean
	go mod tidy

deploy: clean build
	sls deploy --verbose
