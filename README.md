# karmada-operator

Karmada operator for Kubernetes built with Operator SDK and Ansible.This operator is meant to provide a more
Kubernetes-native installation method for Kamrada via a Karmada Custom Resource Definition (CRD). And Karmada-operator 
is an operator which runs as a service on top of Kubernetes.The Karmada-operator servicecan be used to provision and
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
make install
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

## Create a KarmadaDeployment CR

Update the sample KarmadaDeployment CR manifest at `config/samples/operator_v1alpha1_karmadadeployment.yaml` and define the `spec` as the following:

```YAML
apiVersion: operator.karmada.io/v1alpha1
kind: KarmadaDeployment
metadata:
  name: karmadadeployment-sample
spec:
  size: 3
  version: "3.4.9"
```

Create the CR:

```sh
kubectl apply -f config/samples/operator_v1alpha1_karmadadeployment.yaml
```

Ensure that the karmada operator creates the deployment for the sample CR with the correct size:

```console
$ kubectl get deployment
NAME                                    READY   UP-TO-DATE   AVAILABLE   AGE
```

Check the pods and CR status to confirm the status is updated with the karmada pod names:

```console
$ kubectl get pods
NAME                                  READY     STATUS    RESTARTS   AGE
```

```console
$ kubectl get karmadadeployment/karmadadeployment-sample -o yaml
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
