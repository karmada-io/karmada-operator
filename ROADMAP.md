Karmada operator for Kubernetes built with Operator SDK and Ansible.This operator is meant to provide a more Kubernetes-native installation method for Kamrada via a Karmada Custom Resource Definition (CRD).

**Karmada Operator CI**
- [ ]  Create VM by kubevirt
- [ ]  CI test matrix of karmada clsuter lifecyle
- [ ] unit test

**Karmada Control Plane lifecycle**（detailed logic: [reconcile_karmada_controlplane](https://github.com/vivo/karmada-operator/blob/main/roles/karmadadeployment/tasks/reconcile_karmada_controlplane.yaml)）
- [X] install karmada control plane
- [ ]  upgrade karmada control plane
- [ ]  scale karmada control plane
- [X] create karmada-apiserver kubeconfig  on host cluster([karmada-apiserver](https://github.com/vivo/karmada-operator/blob/main/roles/karmadadeployment/tasks/create_karmada_conf.yaml), [k8s-job-install](https://github.com/vivo/karmada-operator/blob/main/roles/karmadadeployment/templates/karmada_config_job.yaml))

**Etcd cluster lifecycle**(detailed logic: [etcd-statefulset](https://github.com/vivo/karmada-operator/blob/main/roles/karmadadeployment/templates/etcd-ss.yaml), [etcd ansible plugins](https://github.com/vivo/karmada-operator/blob/main/roles/karmadadeployment/library/etcd_member.py))
- [X] install etcd cluster
- [ ]  backup etcd cluster
- [ ]  restroe etcd cluster
- [ ]  upgrade etcd
- [X] join etcd member
- [X] remove etcd member 

**Member k8s cluster lifecycle**
- [X] join member cluster([support push and pull mode](https://github.com/vivo/karmada-operator/tree/main/roles/add_karmada_member/tasks))
- [ ] unjoin member cluster([ansible inventory plugins need support](https://github.com/vivo/karmada-operator/blob/main/plugins/inventory/karmada.py#L217), [delete-member-logic](https://github.com/vivo/karmada-operator/tree/main/roles/del_karmada_member/tasks))

**Use binary install karmada**
- [ ] [ansible Karmada plugins Parse the karmada deployment and login with ssh](https://github.com/vivo/karmada-operator/blob/main/plugins/inventory/karmada.py)
- [ ] etcd clsuter

**Karmada cluster Addons**

