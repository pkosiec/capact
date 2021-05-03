#!/usr/bin/env bash

echo "Setting up CI environmental variables for non GCP..."
export NAME="dev"

# LOAD_BALANCER_EXTERNAL_IP is a reserved IP in "External IP addresses" on GCP. It needs to be in the same region.
# Remember when changing LOAD_BALANCER_EXTERNAL_IP to update record A in the Cloud DNS for gateway
cat <<EOT >> "$GITHUB_ENV"
GO_VERSION=^1.16.3
SKIP_DEPS_INSTALLATION=false
CERT_MAX_AGE=85
CERT_NUMBER_TO_BACKUP=1
CERT_SERVICE_NAMESPACE=capact-system
EOT


if [ "${GITHUB_EVENT_NAME}" = "pull_request" ]
then
  PR_NUMBER=$(echo "$GITHUB_REF" | awk 'BEGIN { FS = "/" } ; { print $3 }')
  echo "DOCKER_TAG=PR-${PR_NUMBER}" >> "$GITHUB_ENV"
  echo "DOCKER_REPOSITORY=ghcr.io/project-voltron/go-voltron/pr" >> "$GITHUB_ENV"
else
  echo "DOCKER_TAG=${GITHUB_SHA:0:7}" >> "$GITHUB_ENV"
  echo "DOCKER_REPOSITORY=gcr.io/project-voltron/go-voltron" >> "$GITHUB_ENV"
fi

function returnInfraMatrixIfNeeded() {
  while read -r file; do
    if [[ $file == hack/images/* ]]; then
      # TODO: jinja2 is a Capact Action. Move it to a separate directory or create a new repo for it
      echo 'INFRAS=name=matrix::{"include":[{"INFRA":"json-go-gen"},{"INFRA":"graphql-schema-linter"},{"INFRA":"jinja2"}]}'
      break
    fi
  done <<< "$(gitChanges)"
}

function gitChanges() {
  local DIFF
  # See https://github.community/t/check-pushed-file-changes-with-git-diff-tree-in-github-actions/17220/10
  if [ "$GITHUB_BASE_REF" ]; then
    # Pull Request
    git fetch origin "$GITHUB_BASE_REF" --depth=1
    DIFF=$( git diff --name-only origin/"$GITHUB_BASE_REF" "$GITHUB_SHA" )
  else
    # Push
    DIFF=$( git diff --name-only HEAD^ HEAD )
  fi

  echo "$DIFF"
}

# TODO: Read components to build in automated way, e.g. from directory structure
cat <<EOT >>"$GITHUB_ENV"
APPS=name=matrix::{"include":[{"APP":"gateway"},{"APP":"k8s-engine"},{"APP":"och-js"},{"APP":"argo-runner"},{"APP":"helm-runner"},{"APP":"cloudsql-runner"},{"APP":"populator"},{"APP":"terraform-runner"},{"APP":"argo-actions"}]}
TESTS=name=matrix::{"include":[{"TEST":"e2e"}]}
TOOLS=name=matrix::{"include":[{"TOOL":"ocftool"}]}
$(returnInfraMatrixIfNeeded)
EOT
