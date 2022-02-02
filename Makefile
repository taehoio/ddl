.PHONY: install-dependencies
install-dependencies:
	@go install \
		github.com/bufbuild/buf/cmd/buf@v1.0.0-rc10 \
		github.com/bufbuild/buf/cmd/protoc-gen-buf-breaking@v1.0.0-rc10 \
		github.com/bufbuild/buf/cmd/protoc-gen-buf-lint@v1.0.0-rc10
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.27.1
	@go install github.com/taehoio/protoc-gen-go-ddl@3e26da320ec84fc927d31fef547df0658227aa66

.PHONY: lint
lint: install-dependencies
	buf lint

.PHONY: generate
generate: install-dependencies
	buf generate
	make mock

.PHONY: clean
clean:
	rm -rf gen

.PHONY: diff
diff:
	git diff --exit-code
	if [ -n "$(git status --porcelain)" ]; then git status; exit 1; else exit 0; fi

.PHONY: mock
mock:
	@go install golang.org/x/tools/cmd/stringer@latest
	@go install github.com/golang/mock/mockgen@v1.6.0
	go generate ./...
