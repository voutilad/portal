KEYFILE = "./.secrets/portal-test.json"

portal: cmd/portal/main.go
	go build github.com/voutilad/portal/cmd/portal

clean:
	rm portal

run: portal
	@GOOGLE_APPLICATION_CREDENTIALS=$(KEYFILE) ./portal

.PHONY: clean run
