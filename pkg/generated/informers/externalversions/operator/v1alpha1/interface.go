// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	internalinterfaces "github.com/karmada-io/karmada-operator/pkg/generated/informers/externalversions/internalinterfaces"
)

// Interface provides access to all the informers in this group version.
type Interface interface {
	// Karmadas returns a KarmadaInformer.
	Karmadas() KarmadaInformer
	// KarmadaOnborads returns a KarmadaOnboradInformer.
	KarmadaOnborads() KarmadaOnboradInformer
}

type version struct {
	factory          internalinterfaces.SharedInformerFactory
	namespace        string
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// New returns a new Interface.
func New(f internalinterfaces.SharedInformerFactory, namespace string, tweakListOptions internalinterfaces.TweakListOptionsFunc) Interface {
	return &version{factory: f, namespace: namespace, tweakListOptions: tweakListOptions}
}

// Karmadas returns a KarmadaInformer.
func (v *version) Karmadas() KarmadaInformer {
	return &karmadaInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}

// KarmadaOnborads returns a KarmadaOnboradInformer.
func (v *version) KarmadaOnborads() KarmadaOnboradInformer {
	return &karmadaOnboradInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}
