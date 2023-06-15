SHELL=/bin/sh

.PHONY: local-build
local-build:
	docker build -t localhost:5001/piper:latest .

.PHONY: init-kind
init-kind:
	@if [ "$( kind get clusters | grep piper)" = "" ]; then sh ./scripts/init-kind.sh; else echo "Kind piper exists, switching context"; fi
	kubectl config set-context kind-piper

.PHONY: deploy
deploy: local-build init-kind
	docker push localhost:5001/piper:latest
	helm upgrade --install piper ./helm-chart -f values.dev.yaml

.PHONY: clean
clean:
	docker stop kind-registry && docker rm kind-registry
	kind delete cluster --name piper