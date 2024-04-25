package kubernetes

import (
	"bytes"
	"context"
	"io"
	"k8s.io/client-go/rest"

	"k8s-delegater/pkg/common"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/azure"
)

type KubernetesClient struct {
	clientset *kubernetes.Clientset
}

func GetRestConfigInCluster() (*rest.Config, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		common.GetLog().Errorf("GetKubeConfigBytes error: %v", err)
		return nil, err
	}

	return config, nil
}

func GetRestConfigFromCredentials(fromEnv bool, credentials *AzureKubernetesCredentials, ctx context.Context) (*rest.Config, error) {
	var azK8SConfig *AzureKubernetesConfig

	if fromEnv {
		azK8SConfig = NewAzureKubernetesConfig(NewAzureKubernetesCredentialsFromEnv())
	} else {
		azK8SConfig = NewAzureKubernetesConfig(credentials)
	}

	clientConfig, err := azK8SConfig.GetClientConfig(ctx)
	if err != nil {
		common.GetLog().Errorf("GetKubernetesClientConfigWithEnv error: %v", err)
		return nil, err
	}

	config, err := clientConfig.ClientConfig()
	if err != nil {
		common.GetLog().Errorf("clientcmd NewDefaultClientConfig error: %v", err)
		return nil, err
	}

	return config, nil
}

func NewKubernetesClientFromRestConfig(restConfig *rest.Config) (*KubernetesClient, error) {
	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		common.GetLog().Errorf("GetKubernetesClientSet error: %v", err)
		return nil, err
	}

	return &KubernetesClient{
		clientset: clientset,
	}, nil
}

func (k *KubernetesClient) GetPodLog(ctx context.Context, namespace string, podName string) (log string, err error) {
	req := k.clientset.CoreV1().Pods(namespace).GetLogs(podName, &v1.PodLogOptions{})
	podLogs, err := req.Stream(ctx)
	if err != nil {
		common.GetLog().Errorf("req Stream error: %v", err)
		return "", err
	}
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		common.GetLog().Errorf("io Copy error: %v", err)
		return "", err
	}
	result := buf.String()
	return result, nil
}

func (k *KubernetesClient) GetPodsByJob(ctx context.Context, namespace string, jobName string) (podList *v1.PodList, err error) {
	pods, err := k.clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{LabelSelector: "job-name=" + jobName})
	if err != nil {
		common.GetLog().Errorf("get pod by job error: %v", err)
		return nil, err
	}
	return pods, nil
}

func (k *KubernetesClient) GetPodsByDeployment(ctx context.Context, namespace string, deploymentName string) (podList *v1.PodList, err error) {
	pods, err := k.clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{LabelSelector: "app=" + deploymentName})
	if err != nil {
		common.GetLog().Errorf("get pod by deployment error: %v", err)
		return nil, err
	}
	return pods, nil
}

func (k *KubernetesClient) DeleteDeployment(ctx context.Context, namespace string, deploymentName string) {
	deletePolicy := metav1.DeletePropagationForeground
	if err := k.clientset.AppsV1().Deployments(namespace).Delete(ctx, deploymentName, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		panic(err)
	}
}
