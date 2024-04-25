package helm

import (
	"context"
	"io/ioutil"
	"k8s-delegater/pkg/common/request"
	"k8s-delegater/pkg/dao/keyvault"
	"strings"

	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"sigs.k8s.io/yaml"

	"k8s-delegater/pkg/common"
)

func OverwriteEnvValues(envMap map[string]interface{}, flinkJobRequest request.FlinkJobRequest) (map[string]interface{}, error) {
	envMap["artifactId"] = flinkJobRequest.ArtifactId
	envMap["buildId"] = flinkJobRequest.BuildId
	envMap["dwLayer"] = flinkJobRequest.DWLayer
	envMap["application"] = flinkJobRequest.Application
	envMap["clusterId"] = flinkJobRequest.ClusterId

	kv, err := keyvault.NewKeyVaultClient()
	if err != nil {
		common.GetLog().Errorf("keyvault NewKeyVaultClient error: %v", err)
		return nil, err
	}

	accountKey, err := kv.GetSecret(context.Background(), "XPayADLSGen2AccountKey")
	if err != nil {
		common.GetLog().Errorf("keyvault GetSecret error: %v", err)
		return envMap, err
	}
	envMap["adlsGen2AccountKey"] = accountKey
	return envMap, nil
}

func BuildFlinkParameters(flinkParameters map[string]string) string {
	var builder strings.Builder
	for key, value := range flinkParameters {
		builder.WriteString(" -D")
		builder.WriteString(key)
		builder.WriteString("=")
		builder.WriteString(value)
	}
	return builder.String()
}

func OverwriteValuesYaml(chart *chart.Chart, flinkJobRequest request.FlinkJobRequest) (*chart.Chart, error) {
	chart.Metadata.Name = flinkJobRequest.ReleaseName

	envMap, err := OverwriteEnvValues(chart.Values["env"].(map[string]interface{}), flinkJobRequest)
	if err != nil {
		common.GetLog().Errorf("keyvault GetSecret error: %v", err)
		return chart, err
	}
	chart.Values["env"] = envMap

	chart.Values["nameOverride"] = "xpay-" + flinkJobRequest.ReleaseName
	chart.Values["fullnameOverride"] = "xpay-" + flinkJobRequest.ReleaseName

	flinkParameters := BuildFlinkParameters(flinkJobRequest.FlinkParameters)
	chart.Values["command"].([]interface{})[2] = strings.Replace(chart.Values["command"].([]interface{})[2].(string), "placeholder", flinkParameters, 1)

	return chart, nil
}

func GetHelmChart(chartPath string) (*chart.Chart, error) {
	chart, err := loader.Load(chartPath)
	if err != nil {
		common.GetLog().Errorf("loader Load err %v", err)
		return nil, err
	}

	chart.Values = make(map[string]interface{})

	yamlFile, err := ioutil.ReadFile(chartPath + "/" + common.GetValuesFile())
	if err != nil {
		common.GetLog().Errorf("ioutil ReadFile err %v ", err)
		return nil, err
	}
	if err := yaml.Unmarshal(yamlFile, &chart.Values); err != nil {
		common.GetLog().Errorf("yaml Unmarshal %v ", err)
		return nil, err
	}

	return chart, nil
}
