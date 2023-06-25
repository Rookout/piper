SHELL := /bin/sh

.PHONY: ngrok
ngrok:
	ngrok http 80

.PHONY: local-build
local-build:
	DOCKER_BUILDKIT=1 docker build -t localhost:5001/piper:latest .

.PHONY: init-kind
init-kind:
ifeq ("$(kind get clusters -q | grep piper)", "")
	sh ./scripts/init-kind.sh
else
	echo "Kind piper exists, skipping cluster installation"
endif
	kubectl config set-context kind-piper

.PHONY: init-nginx
init-nginx: init-kind
ifeq ("$(kubectl get pods -n ingress-nginx | grep nginx)", "")
	sh ./scripts/init-nginx.sh
else
	echo "Nginx controller exists, skipping installation"
endif

.PHONY: init-argo-workflows
init-argo-workflows: init-kind
ifeq ("$(helm list -n workflows | grep argo-workflow)", "")
	sh ./scripts/init-argo-workflows.sh
else
	echo "Workflows release exists, skipping installation"
endif

.PHONY: init-piper
init-piper: init-kind local-build
ifeq ("$(helm list | grep piper)", "")
	helm upgrade --install piper ./helm-chart -f values.dev.yaml
else
	echo "Workflows release exists, skipping installation"
endif

.PHONY: deploy
deploy: init-kind init-nginx init-argo-workflows local-build init-piper
	docker push localhost:5001/piper:latest

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