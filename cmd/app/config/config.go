/*
Copyright 2022 The Karmada operator Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package config

import (
	clientset "k8s.io/client-go/kubernetes"
	rest "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
	componentbaseconfig "k8s.io/component-base/config"

	crdclientset "github.com/daocloud/karmada-operator/pkg/generated/clientset/versioned"
	helminstaller "github.com/daocloud/karmada-operator/pkg/installer/helm"
)

type Config struct {
	Client        *clientset.Clientset
	Kubeconfig    *rest.Config
	CRDClient     *crdclientset.Clientset
	EventRecorder record.EventRecorder

	ChartResource  *helminstaller.ChartResource
	LeaderElection componentbaseconfig.LeaderElectionConfiguration
}

type completedConfig struct {
	*Config
}

// CompletedConfig same as Config, just to swap private object.
type CompletedConfig struct {
	// Embed a private pointer that cannot be instantiated outside of this package.
	*completedConfig
}

// Complete fills in any fields not set that are required to have valid data. It's mutating the receiver.
func (c *Config) Complete() *CompletedConfig {
	cc := completedConfig{c}

	// TODO:

	return &CompletedConfig{&cc}
}
