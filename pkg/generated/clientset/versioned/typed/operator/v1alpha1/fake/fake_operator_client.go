// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1alpha1 "github.com/karmada-io/karmada-operator/pkg/generated/clientset/versioned/typed/operator/v1alpha1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeOperatorV1alpha1 struct {
	*testing.Fake
}

func (c *FakeOperatorV1alpha1) Karmadas(namespace string) v1alpha1.KarmadaInterface {
	return &FakeKarmadas{c, namespace}
}

func (c *FakeOperatorV1alpha1) KarmadaOnborads(namespace string) v1alpha1.KarmadaOnboradInterface {
	return &FakeKarmadaOnborads{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeOperatorV1alpha1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
