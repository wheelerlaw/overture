version=$(shell git describe --tags)

octodns: main.go
	go get
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-s -w -X main.version=$(version)" -o octodns main.go

zip: octodns octodns-linux-amd64.zip

octodns-linux-amd64.zip: octodns LICENSE README.md config.sample.json
	zip octodns-linux-amd64.zip octodns config.sample.json LICENSE README.md

clean:
	rm octodns
	rm octodns-linux-amd64.zip

.PHONY: zip clean