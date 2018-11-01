zip: octodns octodns-linux-amd64.zip

octodns: main.go
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-s -w -X main.version=something" -o octodns main.go

octodns-linux-amd64.zip: octodns LICENSE README.md config.sample.json
	zip octodns-linux-amd64.zip octodns config.sample.json LICENSE README.md

clean:
	rm octodns
	rm octodns-linux-amd64.zip

.PHONY: zip clean