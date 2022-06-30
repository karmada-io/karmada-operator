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
```YAML
apiVersion: v1
stringData:
  kubeconfig: |-
    apiVersion: v1
    clusters:
    - cluster:
        certificate-authority-data: XXXXXXX
        server: https://145.40.9.21:6443
      name: cluster.local
    contexts:
    - context:
        cluster: cluster.local
        user: kubernetes-admin
      name: kubernetes-admin@cluster.local
    current-context: kubernetes-admin@cluster.local
    kind: Config
    preferences: {}
    users:
    - name: kubernetes-admin
      user:
        client-certificate-data: XXXXXX
        client-key-data: XXXXXXX
kind: Secret
metadata:
  name: member1-config
  namespace: karmada-system
```
Create the secret: 

```sh
kubectl apply -f secret.yaml
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
      kubeConfigSecretName: "member1-config"
    - name: "member2"
      syncMode: "pull"
      kubeConfigSecretName: "member2-config"
```

Create the CR:

```sh
kubectl apply -f config/samples/operator_v1alpha1_karmadadeployment.yaml
```

Ensure that the karmada operator creates the deployment for the sample CR with the correct size:

```console
$kubectl get deploy -n karmada-system
NAME                                 READY   STATUS     RESTARTS       AGE
karmada-aggregated-apiserver          1/1     1            1           2m50s
karmada-apiserver                     1/1     1            1           2m53s
karmada-controller-manager            1/1     1            1           2m14s
karmada-descheduler                   1/1     1            1           2m11s
karmada-kube-controller-manager       1/1     1            1           2m52s
karmada-scheduler                     1/1     1            1           2m13s
karmada-scheduler-estimator-member1   2/2     2            2           108s
karmada-webhook                       1/1     1            1           2m10s

```

Check the pods and CR status to confirm the status is updated with the karmada pod names:

```console
$ kubectl get pods -n karmada-system
NAME                                                  READY   STATUS    RESTARTS   AGE
karmada-aggregated-apiserver-74c4bd9976-c7r2b          1/1     Running   0          96s
karmada-apiserver-65f5fd7fbf-q4hbz                     1/1     Running   0          99s
karmada-controller-manager-7b6576dcff-mt4t2            1/1     Running   0          60s
karmada-descheduler-65c75d4448-fjw9b                   1/1     Running   0          57s
karmada-kube-controller-manager-7bfff8589f-tgcfh       1/1     Running   0          98s
karmada-scheduler-5f7796487b-7wnbr                     1/1     Running   0          59s
karmada-scheduler-estimator-member1-6f48c5df7b-56xdj   1/1     Running   0          34s
karmada-scheduler-estimator-member1-6f48c5df7b-pnlq2   1/1     Running   0          34s
karmada-webhook-596df4d86c-hb8fg                       1/1     Running   0          56s
karmadadeployment-sample-etcd-0                        1/1     Running   0          2m29s
karmadadeployment-sample-etcd-1                        1/1     Running   0          2m26s
karmadadeployment-sample-etcd-2                        1/1     Running   0          2m23s
```

```console
$ kubectl get karmadadeployment/karmadadeployment-sample -o yaml
```

```console
$ kubectl get cluster --kubeconfig /etc/karmada/kubeconfig
NAME      VERSION   MODE   READY   AGE
member1   v1.21.6   Push   True    2min
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
