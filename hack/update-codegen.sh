#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail
set -ex

# For all commands, the working directory is the parent directory(repo root).
REPO_ROOT=$(git rev-parse --show-toplevel)
cd "${REPO_ROOT}"

GO111MODULE=on go install k8s.io/code-generator/cmd/deepcopy-gen
GO111MODULE=on go install k8s.io/code-generator/cmd/register-gen
GO111MODULE=on go install k8s.io/code-generator/cmd/conversion-gen
GO111MODULE=on go install k8s.io/code-generator/cmd/client-gen
GO111MODULE=on go install k8s.io/code-generator/cmd/lister-gen
GO111MODULE=on go install k8s.io/code-generator/cmd/informer-gen

echo "Generating with deepcopy-gen"
deepcopy-gen \
  --go-header-file "${REPO_ROOT}/hack/boilerplate/boilerplate.go.txt" \
  --input-dirs=github.com/karmada-io/karmada-operator/pkg/apis/operator/v1alpha1 \
  --output-package=github.com/karmada-io/karmada-operator/pkg/apis/operator/v1alpha1 \
  --output-file-base=zz_generated.deepcopy

echo "Generating with register-gen"
register-gen \
  --go-header-file hack/boilerplate/boilerplate.go.txt \
  --input-dirs=github.com/karmada-io/karmada-operator/pkg/apis/operator/v1alpha1 \
  --output-package=github.com/karmada-io/karmada-operator/pkg/apis/operator/v1alpha1 \
  --output-file-base=zz_generated.register

echo "Generating with conversion-gen"
conversion-gen \
  --go-header-file hack/boilerplate/boilerplate.go.txt \
  --input-dirs=github.com/karmada-io/karmada-operator/pkg/apis/operator/v1alpha1 \
  --output-package=github.com/karmada-io/karmada-operator/pkg/apis/operator/v1alpha1 \
  --output-file-base=zz_generated.conversion

echo "Generating with client-gen"
client-gen \
  --go-header-file hack/boilerplate/boilerplate.go.txt \
  --input-base="" \
  --input=github.com/karmada-io/karmada-operator/pkg/apis/operator/v1alpha1 \
  --output-package=github.com/karmada-io/karmada-operator/pkg/generated/clientset \
  --clientset-name=versioned

echo "Generating with lister-gen"
lister-gen \
  --go-header-file hack/boilerplate/boilerplate.go.txt \
  --input-dirs=github.com/karmada-io/karmada-operator/pkg/apis/operator/v1alpha1 \
  --output-package=github.com/karmada-io/karmada-operator/pkg/generated/listers

echo "Generating with informer-gen"
informer-gen \
  --go-header-file "hack/boilerplate/boilerplate.go.txt" \
  --input-dirs "github.com/karmada-io/karmada-operator/pkg/apis/operator/v1alpha1" \
  --versioned-clientset-package github.com/karmada-io/karmada-operator/pkg/generated/clientset/versioned \
  --listers-package github.com/karmada-io/karmada-operator/pkg/generated/listers \
  --output-package "github.com/karmada-io/karmada-operator/pkg/generated/informers"
