TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
WEBSITE_REPO=github.com/hashicorp/terraform-website
PKG_NAME=nsxt
GIT_COMMIT=$$(git rev-list -1 HEAD)
BUILD_PATH=$$(go env GOPATH)

default: build

tools:
	GO111MODULE=on go install -mod=mod github.com/golangci/golangci-lint/cmd/golangci-lint@v1.45
	GO111MODULE=on go install -mod=mod github.com/katbyte/terrafmt

build: fmtcheck
	go install -ldflags "-X github.com/vmware/terraform-provider-nsxt/nsxt.GitCommit=$(GIT_COMMIT)"

build-coverage:
	go build -cover -ldflags "-X github.com/vmware/terraform-provider-nsxt/nsxt.GitCommit=$(GIT_COMMIT)" -o $(BUILD_PATH)/bin

build-debug:
	go build -gcflags="all=-N -l"  -ldflags "-X github.com/vmware/terraform-provider-nsxt/nsxt.GitCommit=$(GIT_COMMIT)" -o $(BUILD_PATH)/bin

test: fmtcheck
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc: fmtcheck
	GO111MODULE=on TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 360m

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
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

website-lint:
	@echo "==> Checking HCL formatting for docs..."
	@terrafmt diff ./docs --check --pattern '*.md' --quiet || (echo; \
	    echo "Unexpected differences in HCL formatting for docs."; \
	    echo "To see the full differences, run: terrafmt diff ./docs --pattern '*.md'"; \
	    echo "To automatically fix the formatting, run 'make docs-lint-fix' and commit the changes."; \
	    exit 1)

docs-lint-fix:
	@echo "==> Applying HCL formatting for docs..."
	@terrafmt fmt ./docs --pattern '*.md'

docs-list-category:
	@find . -name *.md | xargs grep subcategory | awk  -F '"' '{print $$2}' | sort | uniq

.PHONY: build test testacc vet fmt fmtcheck errcheck test-compile docs-lint docs-lint-fix tools

api-wrapper:
	@echo "==> Generating API wrappers..."
	/usr/bin/python3 $(CURDIR)/tools/api-wrapper-generator.py \
		--api_list $(CURDIR)/api/api_list.yaml \
		--api_template $(CURDIR)/api/api_templates.yaml \
		--api_file_template $(CURDIR)/api/api_file_template.yaml \
		--utl_file_template $(CURDIR)/api/utl_file_template.yaml \
		--out_dir $(CURDIR)/api
