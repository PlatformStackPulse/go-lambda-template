#!/usr/bin/env bash

set -euo pipefail

IMAGE_TAG=${1:-latest}
TERRAFORM_DIR=${2:-deploy/terraform}
LOCAL_IMAGE=${3:-go-lambda-template:${IMAGE_TAG}}

if ! command -v docker >/dev/null 2>&1; then
  echo "docker is required"
  exit 1
fi

if ! command -v aws >/dev/null 2>&1; then
  echo "aws CLI is required"
  exit 1
fi

REPOSITORY_URL=$(terraform -chdir="${TERRAFORM_DIR}" output -raw ecr_repository_url)
REGISTRY_HOST=${REPOSITORY_URL%%/*}
AWS_REGION=$(echo "${REGISTRY_HOST}" | cut -d'.' -f4)

aws ecr get-login-password --region "${AWS_REGION}" | docker login --username AWS --password-stdin "${REGISTRY_HOST}"
docker tag "${LOCAL_IMAGE}" "${REPOSITORY_URL}:${IMAGE_TAG}"
docker push "${REPOSITORY_URL}:${IMAGE_TAG}"

echo "Pushed ${REPOSITORY_URL}:${IMAGE_TAG}"
