PROJECT		= "github.com/kula/vault-plugin-secrets-backblazeb2"
GOFILES		= $(shell find . -name "*.go")

default: vault-plugin-secrets-backblazeb2

vault-plugin-secrets-backblazeb2: $(GOFILES)
	go build ./cmd/vault-plugin-secrets-backblazeb2

clean:
	rm -f vault-plugin-secrets-backblazeb2

test:
	/bin/bash test/test.sh

.PHONY: default clean test
