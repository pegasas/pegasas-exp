package flink

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "k8s.io/client-go/plugin/pkg/client/auth/azure"

	"k8s-delegater/pkg/common"
	"k8s-delegater/pkg/common/request"
	"k8s-delegater/pkg/common/response"
	"k8s-delegater/pkg/dao/kubernetes"
	"k8s-delegater/pkg/service/helm"
)

const (
	JOB_NAMESPACE   = "job"
	FLINK_NAMESPACE = "flink"
	CHART_PATH      = "deployment/flink-deployer"
)

func CreateFlinkDeploymentJob(ginContext *gin.Context) {

	var flinkJobRequest request.FlinkJobRequest
	err := ginContext.BindJSON(&flinkJobRequest)
	if err != nil {
		common.GetLog().Errorf("flinkJobRequest bindJson error: %v", err)
		panic(response.Fail(http.StatusInternalServerError, err.Error(), nil))
	}

	// CheckFlinkJobRequest(flinkJobRequest)

	restConfig, err := kubernetes.GetRestConfigInCluster()
	if err != nil {
		common.GetLog().Errorf("kubernetes GetRestConfigInCluster error: %v", err)
		panic(response.Fail(http.StatusInternalServerError, err.Error(), nil))
	}

	actionConfig, err := helm.GetActionConfigFromRestConfig(restConfig, JOB_NAMESPACE)
	if err != nil {
		common.GetLog().Errorf("helm GetActionConfig error: %v", err)
		panic(response.Fail(http.StatusInternalServerError, err.Error(), nil))
	}

	helmClient := helm.NewHelmFromActionConfig(actionConfig)

	releases, err := helmClient.HelmList()
	if err != nil {
		common.GetLog().Errorf("helm HelmList error: %v", err)
		panic(response.Fail(http.StatusInternalServerError, err.Error(), nil))
	}

	releaseName := flinkJobRequest.ReleaseName
	hasRelease := false
	for _, release := range releases {
		if release.Name == releaseName {
			hasRelease = true
			break
		}
	}

	if hasRelease {
		err = helmClient.HelmUninstall(releaseName)
		if err != nil {
			common.GetLog().Errorf("helm HelmUninstall error: %v", err)
			panic(response.Fail(http.StatusInternalServerError, err.Error(), nil))
		}
	}

	chart, err := helm.GetHelmChart(CHART_PATH)
	if err != nil {
		common.GetLog().Errorf("helm GetHelmChart error: %v", err)
		panic(response.Fail(http.StatusInternalServerError, err.Error(), nil))
	}

	chart, err = helm.OverwriteValuesYaml(chart, flinkJobRequest)
	if err != nil {
		common.GetLog().Errorf("helm OverwriteEnvADLSGen2AccountKey error: %v", err)
		panic(response.Fail(http.StatusInternalServerError, err.Error(), nil))
	}

	err = helmClient.HelmInstall(chart, releaseName)
	if err != nil {
		common.GetLog().Errorf("helm HelmInstall error: %v", err)
		panic(response.Fail(http.StatusInternalServerError, err.Error(), nil))
	}

	ginContext.JSON(http.StatusOK, response.Success("success", nil))
}

func GetFlinkDeploymentStatus(ginContext *gin.Context) {
	releaseName := ginContext.Query("releaseName")
	if releaseName == "" {
		common.GetLog().Error("release name is null!")
		panic(response.Fail(http.StatusInternalServerError, "release name is null", nil))
	}

	ctx := context.Background()

	restConfig, err := kubernetes.GetRestConfigInCluster()
	if err != nil {
		common.GetLog().Errorf("kubernetes GetRestConfigInCluster error: %v", err)
		panic(response.Fail(http.StatusInternalServerError, err.Error(), nil))
	}

	k8sClient, err := kubernetes.NewKubernetesClientFromRestConfig(restConfig)
	if err != nil {
		common.GetLog().Errorf("kubernetes NewKubernetesClientFromRestConfig error: %v", err)
		panic(response.Fail(http.StatusInternalServerError, err.Error(), nil))
	}

	pods, err := k8sClient.GetPodsByJob(ctx, JOB_NAMESPACE, "xpay-"+releaseName)
	if err != nil {
		common.GetLog().Errorf("kubernetes GetPodsByJob error: %v", err)
		panic(response.Fail(http.StatusInternalServerError, err.Error(), nil))
	}

	data := make(map[string]string)
	if len(pods.Items) == 0 {
		data["status"] = "UNKNOWN"
	} else {
		data["status"] = string(pods.Items[0].Status.Phase)
	}
	ginContext.JSON(http.StatusOK, response.Success("success", data))
}

func GetFlinkDeploymentLog(ginContext *gin.Context) {
	releaseName := ginContext.Query("releaseName")
	if releaseName == "" {
		common.GetLog().Error("release name is null!")
		panic(response.Fail(http.StatusInternalServerError, "release name is null", nil))
	}

	ctx := context.Background()

	restConfig, err := kubernetes.GetRestConfigInCluster()
	if err != nil {
		common.GetLog().Errorf("kubernetes GetRestConfigInCluster error: %v", err)
		panic(response.Fail(http.StatusInternalServerError, err.Error(), nil))
	}

	k8sClient, err := kubernetes.NewKubernetesClientFromRestConfig(restConfig)
	if err != nil {
		common.GetLog().Errorf("kubernetes NewKubernetesClientFromRestConfig error: %v", err)
		panic(response.Fail(http.StatusInternalServerError, err.Error(), nil))
	}

	pods, err := k8sClient.GetPodsByJob(ctx, JOB_NAMESPACE, "xpay-"+releaseName)
	if err != nil {
		common.GetLog().Errorf("kubernetes GetPodsByJob error: %v", err)
		panic(response.Fail(http.StatusInternalServerError, err.Error(), nil))
	}

	podName := pods.Items[0].GetName()

	logStr, err := k8sClient.GetPodLog(ctx, JOB_NAMESPACE, podName)

	ginContext.JSON(http.StatusOK, response.Success("success", logStr))
}

func GetFlinkJobStatus(ginContext *gin.Context) {
	var flinkJobMeta request.FlinkJobRequest
	err := ginContext.BindJSON(&flinkJobMeta)
	if err != nil {
		// Handle error
	}

	req, err := http.NewRequest(http.MethodGet, common.GetFlinkHistoryServerURL()+"jobs/"+flinkJobMeta.JobId, nil)
	if err != nil {
		common.GetLog().Errorf("http NewRequest error: %v", err)
		panic(response.Fail(http.StatusInternalServerError, err.Error(), nil))
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		common.GetLog().Errorf("http Call error: %v", err)
		panic(response.Fail(http.StatusInternalServerError, err.Error(), nil))
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		common.GetLog().Errorf("ioutil ReadAll error: %v", err)
		panic(response.Fail(http.StatusInternalServerError, err.Error(), nil))
	}

	data := make(map[string]interface{})
	err = json.Unmarshal(resBody, &data)
	if err != nil {
		// Handle error
	}

	if _, ok := data["errors"]; ok {
		ginContext.JSON(http.StatusOK, response.Fail(http.StatusInternalServerError, "get flink job status error", data))
	} else {
		ginContext.JSON(http.StatusOK, response.Success("success", data))
	}
}

func TriggerFlinkJobSavepoint(ginContext *gin.Context) {
	var request request.FlinkRequest
	err := ginContext.BindJSON(&request)
	if err != nil {
		common.GetLog().Errorf("ginContext BindJSON error: %v", err)
		panic(response.Fail(http.StatusInternalServerError, err.Error(), nil))
	}

	bodyMap := map[string]string{"target-directory": common.GetHadoopAzureBlobFileSystemSavepointDirectory() + request.DeploymentName}
	bodyBytes, err := json.Marshal(bodyMap)
	if err != nil {
		common.GetLog().Errorf("json Marshal error: %v", err)
		panic(response.Fail(http.StatusInternalServerError, err.Error(), nil))
	}

	req, err := http.NewRequest(http.MethodPost, "http://"+request.DeploymentName+"-rest.flink.svc.cluster.local:8081/jobs/"+request.JobId+"/savepoints/", bytes.NewBuffer(bodyBytes))
	// req, err := http.NewRequest(http.MethodPost, "http://localhost:8081/jobs/"+request.JobId+"/savepoints/", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		common.GetLog().Errorf("http NewRequest error: %v", err)
		panic(response.Fail(http.StatusInternalServerError, err.Error(), nil))
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		common.GetLog().Errorf("http Call error: %v", err)
		panic(response.Fail(http.StatusInternalServerError, err.Error(), nil))
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		common.GetLog().Errorf("ioutil ReadAll error: %v", err)
		panic(response.Fail(http.StatusInternalServerError, err.Error(), nil))
	}

	data := make(map[string]interface{})
	err = json.Unmarshal(resBody, &data)
	if err != nil {
		common.GetLog().Errorf("json Unmarshal error: %v", err)
		panic(response.Fail(http.StatusInternalServerError, err.Error(), nil))
	}

	ginContext.JSON(http.StatusOK, response.Success("success", data))
}

func GetFlinkJobSavepointTriggerStatus(ginContext *gin.Context) {
	var request request.FlinkRequest
	err := ginContext.BindJSON(&request)
	if err != nil {
		common.GetLog().Errorf("ginContext BindJSON error: %v", err)
		panic(response.Fail(http.StatusInternalServerError, err.Error(), nil))
	}

	req, err := http.NewRequest(http.MethodPost, "http://"+request.DeploymentName+"-rest.flink.svc.cluster.local:8081/jobs/"+request.JobId+"/savepoints/"+request.TriggerId, nil)
	// req, err := http.NewRequest(http.MethodGet, "http://localhost:8081/jobs/"+request.JobId+"/savepoints/"+request.TriggerId, nil)
	if err != nil {
		common.GetLog().Errorf("http NewRequest error: %v", err)
		panic(response.Fail(http.StatusInternalServerError, err.Error(), nil))
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		common.GetLog().Errorf("http Call error: %v", err)
		panic(response.Fail(http.StatusInternalServerError, err.Error(), nil))
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		common.GetLog().Errorf("ioutil ReadAll error: %v", err)
		panic(response.Fail(http.StatusInternalServerError, err.Error(), nil))
	}

	data := make(map[string]interface{})
	err = json.Unmarshal(resBody, &data)
	if err != nil {
		common.GetLog().Errorf("json Unmarshal error: %v", err)
		panic(response.Fail(http.StatusInternalServerError, err.Error(), nil))
	}

	ginContext.JSON(http.StatusOK, response.Success("success", data))
}

func GetFlinkJobID(ginContext *gin.Context) {
	releaseName := ginContext.Query("releaseName")
	req, err := http.NewRequest(http.MethodGet, "http://"+releaseName+".flink.svc.cluster.local:8081/jobs/overview", nil)
	//req, err := http.NewRequest(http.MethodGet, "http://localhost:8081/jobs/overview", nil)
	if err != nil {
		common.GetLog().Errorf("http NewRequest error: %v", err)
		panic(response.Fail(http.StatusInternalServerError, err.Error(), nil))
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		common.GetLog().Errorf("http Call error: %v", err)
		panic(response.Fail(http.StatusInternalServerError, err.Error(), nil))
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		common.GetLog().Errorf("ioutil ReadAll error: %v", err)
		panic(response.Fail(http.StatusInternalServerError, err.Error(), nil))
	}

	data := make(map[string]interface{})
	err = json.Unmarshal(resBody, &data)
	if err != nil {
		// Handle error
	}

	if _, ok := data["errors"]; ok {
		ginContext.JSON(http.StatusOK, response.Fail(http.StatusInternalServerError, "get flink job id error", data))
	} else {
		ginContext.JSON(http.StatusOK, response.Success("success", data))
	}
}

func DeleteFlinkJob(ginContext *gin.Context) {
	deploymentName := ginContext.Query("releaseName")
	ctx := context.Background()

	restConfig, err := kubernetes.GetRestConfigInCluster()
	if err != nil {
		common.GetLog().Errorf("kubernetes GetRestConfigInCluster error: %v", err)
		panic(response.Fail(http.StatusInternalServerError, err.Error(), nil))
	}

	k8sClient, err := kubernetes.NewKubernetesClientFromRestConfig(restConfig)
	if err != nil {
		common.GetLog().Errorf("kubernetes NewKubernetesClientFromRestConfig error: %v", err)
		panic(response.Fail(http.StatusInternalServerError, err.Error(), nil))
	}

	k8sClient.DeleteDeployment(ctx, FLINK_NAMESPACE, deploymentName)

	ginContext.JSON(http.StatusOK, response.Success("delete job success", nil))
}
