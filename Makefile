.PHONY: install-dependencies
install-dependencies:
	@go install \
		github.com/bufbuild/buf/cmd/buf@v1.0.0-rc10 \
		github.com/bufbuild/buf/cmd/protoc-gen-buf-breaking@v1.0.0-rc10 \
		github.com/bufbuild/buf/cmd/protoc-gen-buf-lint@v1.0.0-rc10
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.27.1
	@go install github.com/taehoio/protoc-gen-go-ddl@abb8bccfbb7a7aae0cb2697e05f079a36e7b09e6

.PHONY: lint
lint: install-dependencies
	buf lint

.PHONY: generate
generate: install-dependencies
	buf generate

.PHONY: clean
clean:
	rm -rf gen

.PHONY: diff
diff:
	git diff --exit-code
	if [ -n "$(git status --porcelain)" ]; then git status; exit 1; else exit 0; fi
