KEYFILE = "./.secrets/portal-test.json"

portal: deps cmd/portal/main.go portal.go
	go build github.com/voutilad/portal/cmd/portal

deps:
	go get github.com/voutilad/portal

clean:
	rm -f portal

run: portal
	@GOOGLE_APPLICATION_CREDENTIALS=$(KEYFILE) ./portal

.PHONY: clean deps run
