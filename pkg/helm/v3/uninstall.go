package v3

import (
	"helm.sh/helm/v3/pkg/action"

	"github.com/daocloud/karmada-operator/pkg/helm"
)

func (h *HelmV3) Uninstall(releaseName string, opts helm.UninstallOptions) error {
	cfg, err := newActionConfig(h.kubeconfigPath, h.kubeConfig, h.infoLogFunc(opts.Namespace, releaseName), opts.Namespace, "")
	if err != nil {
		return err
	}

	uninstall := action.NewUninstall(cfg)
	uninstallOptions(opts).configure(uninstall)

	_, err = uninstall.Run(releaseName)
	return err
}

type uninstallOptions helm.UninstallOptions

func (opts uninstallOptions) configure(action *action.Uninstall) {
	action.DisableHooks = opts.DisableHooks
	action.DryRun = opts.DryRun
	action.KeepHistory = opts.KeepHistory
	action.Timeout = opts.Timeout
}
