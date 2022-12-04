# Karmada operator roadmap

The Karmada operator is a new project of the Karmada community designed to manage the lifecycle of Karmada. The following are the core features of the Karmada operator that we will work on. It is not a detailed to-do list but represents the overall development direction. It is hoped that more contributors will have a general understanding of the project and participate in it.

## the first version(v0.0.1)

### feature
- support create Karmada instance on exist Kubernetes (include estimate, deschedule)
- support Karmada instance upgrade
- Karmada certificate manager
- join member cluster (both `Pull` & `Push`)

### others
- the `pr` and `issue` template
- actions
  * code lint
  * build mutil-platform image
- makefile & scirpts
- operator & webhook code framework
- build chart package
- sample files
- documents
  * readme.md
  * chart readme.md
  * chart NOTES.txt
  * record demo video
  * update Karmada website
- actions
  * sync Karmada CRD files
  * image release
  * chart-lint
  * chart release
  * unit test
- e2e framework

## the second version

### feature

- ETCD HA
- integrate Karmada sub command
- expose Karmada apiserver with MetalLB
- Karmada components health check
- migrate Karmada cluster
- ETCD data backup & restore
- support metrics

### others

- e2e cases
- actions:
  * e2e test
