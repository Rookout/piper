#!/bin/sh
set -o errexit

# 7. Install argo workflows
helm repo add argo https://argoproj.github.io/argo-helm
helm upgrade --install argo-workflow argo/argo-workflows -n workflows --create-namespace -f workflows.values.yaml