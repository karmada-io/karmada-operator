package main

import (
	"os"

	apiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/component-base/cli"
	_ "k8s.io/component-base/logs/json/register"          // for JSON log format registration
	_ "k8s.io/component-base/metrics/prometheus/clientgo" // load all the prometheus client-go plugin
	_ "k8s.io/component-base/metrics/prometheus/version"  // for version metric registration

	"github.com/karmada-io/karmada-operator/cmd/operator/app"
)

func main() {
	ctx := apiserver.SetupSignalContext()
	command := app.NewOperatorCommand(ctx)
	code := cli.Run(command)
	os.Exit(code)
}
