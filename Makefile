TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
TARGETS=darwin linux windows

default: build

build: fmt
	go install

test: fmt
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc: fmt
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@echo "==> Checking that code complies with gofmt requirements..." ; \
	files=$$(find . -name '*.go' ) ; \
	gofmt_files=`gofmt -l $$files`; \
	if [ -n "$$gofmt_files" ]; then \
		echo 'gofmt needs running on the following files:'; \
		echo "$$gofmt_files"; \
		echo "You can use the command: \`make fmt\` to reformat code."; \
		exit 1; \
	fi

vendor-status:
	@govendor status

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./aws"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

targets: $(TARGETS)

$(TARGETS):
	GOOS=$@ GOARCH=amd64 go build -o "dist/$@/terraform-provider-tfstate_${TRAVIS_TAG}_x4"
	zip -j dist/terraform-provider-tfstate_${TRAVIS_TAG}_$@_amd64.zip dist/$@/terraform-provider-tfstate_${TRAVIS_TAG}_x4

.PHONY: build test testacc vet fmt errcheck vendor-status test-compile targets $(TARGETS)

