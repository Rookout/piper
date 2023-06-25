SHELL := /bin/sh
CLUSTER_DEPLOYED := $(shell kind get clusters -q | grep piper)

.PHONY: ngrok
ngrok:
	ngrok http 80

.PHONY: local-build
local-build:
	DOCKER_BUILDKIT=1 docker build -t localhost:5001/piper:latest .

.PHONY: local-push
local-push:
	docker push localhost:5001/piper:latest

.PHONY: init-kind
init-kind:
ifndef CLUSTER_DEPLOYED
	sh ./scripts/init-kind.sh
else
	$(info Kind piper cluster exists, skipping cluster installation)
endif
	kubectl config set-context kind-piper

.PHONY: init-nginx
init-nginx: init-kind
	sh ./scripts/init-nginx.sh

.PHONY: init-argo-workflows
init-argo-workflows: init-kind
	sh ./scripts/init-argo-workflows.sh

.PHONY: init-piper
init-piper: init-kind local-build
	sh ./scripts/init-piper.sh

.PHONY: deploy
deploy: init-kind init-nginx init-argo-workflows local-build local-push init-piper

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
	helm template ./helm-chart --debug > _lint.yaml
	helm-docs

.PHONY: test
test:
	go test -short ./pkg/...

$(GOPATH)/bin/golangci-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b `go env GOPATH`/bin v1.52.2

.PHONY: lint
lint: $(GOPATH)/bin/golangci-lint
	$(GOPATH)/bin/golangci-lint run --fix --verbose