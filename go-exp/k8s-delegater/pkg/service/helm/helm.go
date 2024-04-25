package helm

import (
	"fmt"
	"os"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/release"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"k8s-delegater/pkg/common"
)

type HelmClient struct {
	actionConfig *action.Configuration
}

func GetActionConfigFromRestConfig(config *rest.Config, releaseNamespace string) (*action.Configuration, error) {
	actionConfig := new(action.Configuration)

	insecure := true
	cliConfig := genericclioptions.NewConfigFlags(false)
	cliConfig.APIServer = &config.Host
	cliConfig.BearerToken = &config.BearerToken
	cliConfig.Namespace = &releaseNamespace
	cliConfig.Insecure = &insecure
	wrapper := func(*rest.Config) *rest.Config {
		return config
	}
	cliConfig.WithWrapConfigFn(wrapper)
	if err := actionConfig.Init(cliConfig, releaseNamespace, os.Getenv("HELM_DRIVER"), func(format string, v ...interface{}) {
		fmt.Sprintf(format, v)
	}); err != nil {
		common.GetLog().Errorf("actionConfig InitConfig error: %v", err)
		return nil, err
	}
	return actionConfig, nil
}

func GetActionConfigFromClientConfig(clientConfig clientcmd.ClientConfig, releaseNamespace string) (config *action.Configuration, err error) {
	actionConfig := new(action.Configuration)
	restClientGetter := NewRESTClientGetter(releaseNamespace, clientConfig)
	if err := actionConfig.Init(restClientGetter, releaseNamespace, os.Getenv("HELM_DRIVER"), func(format string, v ...interface{}) {
		fmt.Sprintf(format, v)
	}); err != nil {
		common.GetLog().Errorf("actionConfig InitConfig error: %v", err)
		return nil, err
	}
	return actionConfig, nil
}

func NewHelmFromActionConfig(actionConfig *action.Configuration) *HelmClient {
	return &HelmClient{
		actionConfig: actionConfig,
	}
}

func (h *HelmClient) HelmInstall(chart *chart.Chart, releaseName string) error {
	iCli := action.NewInstall(h.actionConfig)
	iCli.ReleaseName = releaseName

	rel, err := iCli.Run(chart, nil)
	if err != nil {
		common.GetLog().Errorf("iCli Run %v ", err)
		return err
	}

	common.GetLog().Info(rel)
	return nil
}

func (h *HelmClient) HelmUninstall(releaseName string) error {
	iCli := action.NewUninstall(h.actionConfig)
	_, err := iCli.Run(releaseName)
	if err != nil {
		common.GetLog().Errorf("iCli Run %v ", err)
		return err
	}
	return nil
}

func (h *HelmClient) HelmList() ([]*release.Release, error) {
	iCli := action.NewList(h.actionConfig)

	releases, err := iCli.Run()
	if err != nil {
		common.GetLog().Errorf("iCli Run %v ", err)
		return nil, err
	}
	return releases, nil
}
