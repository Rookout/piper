#!/bin/sh
set -o errexit

if [ -z "$(helm list | grep piper)" ]; then
  # 8. Install Piper
  helm upgrade --install piper ./helm-chart -f values.dev.yaml
else
  echo "Piper release exists, skipping installation"
fi
