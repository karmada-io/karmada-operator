# karmada-operator

Karmada operator for Kubernetes built with Operator SDK and Ansible.This operator is meant to provide a more
Kubernetes-native installation method for Kamrada via a Karmada Custom Resource Definition (CRD). And Karmada-operator 
is an operator which runs as a service on top of Kubernetes.The Karmada-operator service can be used to provision and
perform initial configuration of karmada.

For provisioning Karmada,the karmada-operator will use the current deployment tool.
- [Kubectl Karmada](https://github.com/karmada-io/karmada/tree/master/cmd/kubectl-karmada)
- [Karmada Charts](https://github.com/karmada-io/karmada/tree/master/charts)
- [binary install karmada]
- [karmada yaml](https://github.com/karmada-io/karmada/tree/master/artifacts/deploy)

## Run as a Deployment inside the cluster
Now that a developer is confident in the operator logic, testing the operator
inside of a pod on a Kubernetes cluster is desired. Running as a pod inside a
Kubernetes cluster is preferred for production use.

To build the `karmada-operator` image and push it to a registry:

```sh
make docker-build docker-push IMG=example.com/karmada-operator:v0.0.1
```

Deploy the karmada-operator:

```sh
make deploy IMG=example.com/karmada-operator:v0.0.1
```

Verify that the karmada-operator is up and running:

```console
$ kubectl get deployment -n karmada-operator-system
NAME                     DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
kamrada-operator       1         1         1            1           1m
```

Uninstall the operator:

```sh
 make undeploy
```

## Create members cluster secret.yaml, example:
show member1 kubeconfig:
```sh
cat ./member1
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: XXXXXXX
    server: https://10.233.66.15:6443
  name: cluster.member1
contexts:
- context:
    cluster: cluster.member1
    user: kubernetes-admin
  name: kubernetes-admin@cluster.member1
current-context: kubernetes-admin@cluster.member1
kind: Config
preferences: {}
users:
- name: kubernetes-admin
  user:
    client-certificate-data: XXXXXXX
    client-key-data: XXXXX
```
Create the secret:
```sh
kubectl create secret generic member1-kubeconfig --from-file=member1-kubeconfig=./member1 -n karmada-system
kubectl create secret generic member2-kubeconfig --from-file=member2-kubeconfig=./member2 -n karmada-system
```

## Create a KarmadaDeployment CR

Update the sample KarmadaDeployment CR manifest at `config/samples/operator_v1alpha1_karmadadeployment.yaml` and define the `spec` as the following:

```YAML
apiVersion: operator.karmada.io/v1alpha1
kind: KarmadaDeployment
metadata:
  name: karmadadeployment-sample
  namespace: karmada-system
spec:
  karmadaRegistry: "swr.ap-southeast-1.myhuaweicloud.com"
  kubeRegistry: "k8s.gcr.io"
  clusterDomain: "cluster.local"
  etcd:
     size: 3
     version: "3.4.9"
     pvc:
       # etcd.pvc.storageClass storageClass name of PVC
       storageClass: "local-path"
       # etcd.pvc.size size of PVC
       size: "1Gi"
     # "pvc" means using volumeClaimTemplates
     # "hostPath" means using hostPath
     storageType: "pvc"
  scheduler:
      size: 1
      version: "v1.2.0"
  webhook:
      size: 1
      version: "v1.2.0"
  controllerManager:
    size: 1
    version: "v1.2.0"
  apiServer:
    size: 1
    version: "v1.21.7"
    ## "LoadBalancer" means using LoadBalancer
    ## "ClusterIP" means using ClusterIP
    ## "NodePort" means using NodePort
    serviceType: "ClusterIP"
    ## if loadBalancerApiserverIp is define,
    ## and ip is not none
    ## we will use external lb vip
    loadBalancerApiserverIp: ""
  aggregatedApiServer:
    size: 1
    version: "v1.2.0"
  kubeControllerManager:
    size: 1
    version: "v1.21.7"
  agent:
    size: 2
    version: "v1.2.0"
  descheduler:
    size: 1
    version: "v1.2.0"
  search:
    size: 1
    version: "v1.2.0"
  schedulerEstimator:
    size: 1
    version: "v1.2.0"
  members:
    - name: "member1"
      syncMode: "push"
      kubeConfigSecretName: "member1-kubeconfig"
    - name: "member2"
      syncMode: "pull"
      kubeConfigSecretName: "member2-kubeconfig"
```

Create the CR:

```sh
kubectl apply -f config/samples/operator_v1alpha1_karmadadeployment.yaml
```

Ensure that the karmada operator creates the deployment for the sample CR with the correct size:

```console
$kubectl get deploy -n karmada-system
NAME                                  READY   UP-TO-DATE   AVAILABLE   AGE
karmada-aggregated-apiserver          1/1     1            1           2m23s
karmada-apiserver                     1/1     1            1           2m25s
karmada-controller-manager            1/1     1            1           1m51s
karmada-descheduler                   1/1     1            1           1m49s
karmada-kube-controller-manager       1/1     1            1           2m24s
karmada-scheduler                     1/1     1            1           1m50s
karmada-scheduler-estimator-member1   1/1     1            1           1m2s
karmada-scheduler-estimator-member2   1/1     1            1           1m8s
karmada-search                        1/1     1            1           2m21s
karmada-webhook                       1/1     1            1           1m48s

```

Check the pods and CR status to confirm the status is updated with the karmada pod names:

```console
$ kubectl get pods -n karmada-system
NAME                                                   READY   STATUS      RESTARTS   AGE
install-kubeconfig-xrltz                               0/1     Completed   0          2m9s
karmada-aggregated-apiserver-74c4bd9976-rhmbl          1/1     Running     0          2m49s
karmada-apiserver-65f5fd7fbf-c5xmv                     1/1     Running     0          2m51s
karmada-controller-manager-7b6576dcff-9qgjg            1/1     Running     0          2m17s
karmada-descheduler-65c75d4448-jf8mh                   1/1     Running     0          2m15s
karmada-kube-controller-manager-7bfff8589f-7nzvz       1/1     Running     0          2m50s
karmada-scheduler-5f7796487b-jnf8v                     1/1     Running     0          2m16s
karmada-scheduler-estimator-member1-6f48c5df7b-kmrgr   1/1     Running     0          38s
karmada-scheduler-estimator-member2-d4b55f54b-8tm9k    1/1     Running     0          34s
karmada-search-647589f58b-sb4m9                        1/1     Running     0          2m47s
karmada-webhook-596df4d86c-j69c9                       1/1     Running     0          2m14s
karmadadeployment-sample-etcd-0                        1/1     Running     0          3m28s
karmadadeployment-sample-etcd-1                        1/1     Running     0          3m27s
karmadadeployment-sample-etcd-2                        1/1     Running     0          3m25s
```

```console
$ kubectl get karmadadeployment/karmadadeployment-sample -o yaml -n karmada-system
apiVersion: operator.karmada.io/v1alpha1
kind: KarmadaDeployment
metadata:
  creationTimestamp: "2022-09-17T14:07:01Z"
  generation: 1
  name: karmadadeployment-sample
  namespace: karmada-system
  resourceVersion: "484749"
  uid: a6ddc1d5-e688-4697-a7b6-9aff167ed176
spec:
  agent:
    size: 2
    version: v1.2.0
  aggregatedApiServer:
    size: 1
    version: v1.2.0
  apiServer:
    loadBalancerApiserverIp: ""
    serviceType: ClusterIP
    size: 1
    version: v1.21.7
  clusterDomain: cluster.local
  controllerManager:
    size: 1
    version: v1.2.0
  descheduler:
    size: 1
    version: v1.2.0
  etcd:
    pvc:
      size: 1Gi
      storageClass: local-path
    size: 3
    storageType: pvc
    version: 3.4.9
  karmadaRegistry: swr.ap-southeast-1.myhuaweicloud.com
  kubeControllerManager:
    size: 1
    version: v1.21.7
  kubeRegistry: k8s.gcr.io
  members:
  - kubeConfigSecretName: member1-kubeconfig
    name: member1
    syncMode: push
  - kubeConfigSecretName: member2-kubeconfig
    name: member2
    syncMode: push
  scheduler:
    size: 1
    version: v1.2.0
  schedulerEstimator:
    size: 1
    version: v1.2.0
  search:
    size: 1
    version: v1.2.0
  webhook:
    size: 1
    version: v1.2.0
status:
  clientPort: 2379
  conditions:
  - lastTransitionTime: "2022-09-17T14:09:17Z"
    message: ""
    reason: ""
    status: "False"
    type: Failure
  - ansibleResult:
      changed: 1
      completion: 2022-09-17T14:11:58.77319
      failures: 0
      ok: 15
      skipped: 38
    lastTransitionTime: "2022-09-17T14:07:01Z"
    message: Awaiting next reconciliation
    reason: Successful
    status: "True"
    type: Running
  - lastTransitionTime: "2022-09-17T14:11:59Z"
    message: Last reconciliation succeeded
    reason: Successful
    status: "True"
    type: Successful
  controlPaused: false
  currentVersion: 3.4.9
  etcdMembers:
    ready: member1,member2
    unready: karmadadeployment-sample-etcd-0,karmadadeployment-sample-etcd-1,karmadadeployment-sample-etcd-2
  etcdPhase: Running
  karmadaPhase: Unknown
  serviceName: karmadadeployment-sample-etcd-client
  size: 3
  targetVersion: 3.4.9
```

```console
$ kubectl get cluster --kubeconfig /etc/karmada/kubeconfig
NAME      VERSION   MODE   READY   AGE
member1   v1.21.6   Push   True    3m3s
member2   v1.21.6   pull   True    3m20s
```

## Cleanup

Run the following to delete all deployed resources:

```sh
kubectl delete -f config/samples/operator_v1alpha1_karmadadeployment.yaml
make undeploy
```

## Viewing the Ansible logs

In order to see the logs from a particular operator you can run:

```sh
kubectl logs deployment/karmada-operator-controller-manager -n karmada-operator-system
```

The logs contain the information about the Ansible run and are useful for
debugging your Ansible tasks. Note that the logs may contain much more
detailed information about the Ansible Operator's internals and its
interactions with Kubernetes as well.

Also, you can set the environment variable `ANSIBLE_DEBUG_LOGS` to `True` to
check the full Ansible result in the logs in order to be able to debug it.

**Example**

In `config/manager/manager.yaml` and `config/default/manager_auth_proxy_patch.yaml`:

```yaml
...
      containers:
      - name: manager
        env:
        - name: ANSIBLE_DEBUG_LOGS
          value: "True"
...
```
