#!/bin/sh
set -o errexit

if [ -z "$(kubectl get pods --all-namespaces | grep ingress-nginx-controller)" ]; then
  # 6. Deploy of nginx ingress controller to the cluster
  kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml && \
  kubectl wait --namespace ingress-nginx \
         --for=condition=complete job/ingress-nginx-admission-create \
         --timeout=180s && \
  kubectl rollout restart deployment ingress-nginx-controller --namespace ingress-nginx && \
  kubectl wait --namespace ingress-nginx \
         --for=condition=ready pod \
         --selector=app.kubernetes.io/component=controller \
         --timeout=360s
else
  echo "Nginx already exists, skipping installation"
fi
