package v3

import (
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/release"

	"github.com/daocloud/karmada-operator/pkg/helm"
)

func (h *HelmV3) UpgradeFromPath(chartPath string, releaseName string, values []byte,
	opts helm.UpgradeOptions) (*helm.Release, error) {

	cfg, err := newActionConfig(h.kubeconfigPath, h.kubeConfig, h.infoLogFunc(opts.Namespace, releaseName), opts.Namespace, "")
	if err != nil {
		return nil, err
	}

	// Load the chart from the given path, this also ensures that
	// all chart dependencies are present
	chartRequested, err := loader.Load(chartPath)
	if err != nil {
		return nil, err
	}

	// Read and set values
	val, err := chartutil.ReadValues(values)
	if err != nil {
		return nil, err
	}

	var res *release.Release
	if opts.Install {
		install := action.NewInstall(cfg)
		installOptions(opts).configure(install, releaseName)
		res, err = install.Run(chartRequested, val.AsMap())
	} else {
		upgrade := action.NewUpgrade(cfg)
		upgradeOptions(opts).configure(upgrade)
		res, err = upgrade.Run(releaseName, chartRequested, val.AsMap())
	}

	if err != nil {
		return nil, err
	}
	return releaseToGenericRelease(res), err
}

type installOptions helm.UpgradeOptions

func (opts installOptions) configure(action *action.Install, releaseName string) {
	action.Namespace = opts.Namespace
	action.ReleaseName = releaseName
	action.Atomic = opts.Atomic
	action.DisableHooks = opts.DisableHooks
	action.DryRun = opts.DryRun
	action.ClientOnly = opts.ClientOnly
	action.Timeout = opts.Timeout
	action.Wait = opts.Wait
	action.SkipCRDs = opts.SkipCRDs
	action.DisableOpenAPIValidation = opts.DisableValidation
}

type upgradeOptions helm.UpgradeOptions

func (opts upgradeOptions) configure(action *action.Upgrade) {
	action.Namespace = opts.Namespace
	action.Atomic = opts.Atomic
	action.DisableHooks = opts.DisableHooks
	action.DryRun = opts.DryRun
	action.Force = opts.Force
	action.MaxHistory = opts.MaxHistory
	action.ResetValues = opts.ResetValues
	action.ReuseValues = opts.ReuseValues
	action.Timeout = opts.Timeout
	action.Wait = opts.Wait
}
