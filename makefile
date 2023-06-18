SHELL=/bin/sh

.PHONY: ngrok
ngrok:
	ngrok http 80

.PHONY: local-build
local-build:
	docker build -t localhost:5001/piper:latest .

.PHONY: init-kind
init-kind:
	@if [[ "$(kind get clusters -q | grep piper)" == "" ]]; then sh ./scripts/init-kind.sh; else echo "Kind piper exists, switching context"; fi
	kubectl config set-context kind-piper

.PHONY: deploy
deploy: local-build init-kind
	docker push localhost:5001/piper:latest
	helm upgrade --install piper ./helm-chart -f values.dev.yaml && kubectl rollout restart deployment piper

.PHONY: restart
restart: local-build
	docker push localhost:5001/piper:latest
	kubectl rollout restart deployment piper

.PHONY: clean
clean:
	docker stop kind-registry && docker rm kind-registry
	kind delete cluster --name piper

.PHONY: helm
helm:
	helm lint ./helm-chart
	helm template ./helm-chart  --debug > _lint.yaml
	helm-docs

.PHONY: test
test:
	go test -short ./pkg/...

$(GOPATH)/bin/golangci-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b `go env GOPATH`/bin v1.52.2