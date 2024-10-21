.PHONY: build-dev
build-dev:
	mkdir -p bin
	go build -tags dev -ldflags "-s -w -X 'github.com/nilpntr/certmetrics-exporter/cmd.version=dev'" -o bin/certmetrics-exporter github.com/nilpntr/certmetrics-exporter

.PHONY: run
run: build-dev
	KUBE_ENV=dev LOG_LEVEL=debug ./bin/certmetrics-exporter

.PHONY: version
version: build-dev
	@./bin/certmetrics-exporter version