GO ?= go
PYTHON ?= python3
PHP ?= $(CURDIR)/containers/php-version/bin/php84
PHP_VERSION ?= 8.4.26
PHP_TARGET = $(basename $(PHP_VERSION))
GOPHP_FILE ?= examples/report.php
GOPHP_TARGET ?= 8.5
# Keep build cache inside the writable workspace, including sandboxed runs.
export GOCACHE := $(CURDIR)/.cache/go-build

.PHONY: build install test vet compare compare-parser compare-all container-build container-run container-check
build:
	$(GO) build -o bin/gophp ./cmd/gophp

# Put gophp on PATH for the current user (override with BINDIR=/usr/local/bin).
BINDIR ?= $(HOME)/.local/bin
install: build
	install -d "$(BINDIR)"
	install -m 0755 bin/gophp "$(BINDIR)/gophp"

test:
	$(GO) test ./...
	$(PYTHON) -B -m unittest discover -s tests -p 'test_*.py'

vet:
	$(GO) vet ./...

# Run local expectations first: success fixtures must actually succeed.
compare: test build
	$(PYTHON) skills/php-go-interpreter/scripts/compare_php.py --php "$(PHP)" --php-version "$(PHP_VERSION)" --fixtures tests/fixtures --interpreter ./bin/gophp --php-version=$(PHP_TARGET) '{file}'

# Check parser diagnostics separately when investigating syntax errors.
compare-parser: test build
	$(PYTHON) skills/php-go-interpreter/scripts/compare_php.py --php "$(PHP)" --php-version "$(PHP_VERSION)" --fixtures tests/fixtures/parser --interpreter ./bin/gophp --php-version=$(PHP_TARGET) '{file}'

# Compare all fixtures against all three official PHP versions.
compare-all: test build
	./containers/php-version/compare-all.sh tests/fixtures ./bin/gophp

# Run the Go interpreter itself in a one-shot Podman container.
container-build:
	./containers/gophp/build.sh

container-run:
	./containers/gophp/run.sh --php-version=$(GOPHP_TARGET) "$(GOPHP_FILE)"

container-check:
	$(PYTHON) -B containers/gophp/check.py
