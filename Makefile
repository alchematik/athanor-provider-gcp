.PHONY: install/athanor
install/athanor:
	cd ../athanor && go install ./cmd/athanor && cd -

.PHONY: build/translator/go
build/translator/go:
	cd ../athanor-go/cmd/translator && go build -o ../../../athanor-provider-gcp/build/translator/go/v0.0.1/translator && cd -

.PHONY: generate
generate: install/athanor build/translator/go
	athanor provider generate config.json

.PHONY: build/provider
build/provider:
	go build -o build/provider/gcp/v0.0.1/provider ./cmd/provider

.PHONY: blueprint/reconcile
blueprint/reconcile: install/athanor build/translator/go build/provider
	athanor blueprint reconcile ./example/config.json

.PHONY: diff/show
diff/show: install/athanor
	athanor diff show --debug ./example/config.json 

.PHONY: diff/reconcile
diff/reconcile: install/athanor
	athanor diff reconcile --debug ./example/config.json

.PHONY: deps/install
deps/install: install/athanor
	athanor deps install --debug ./example/config.json

