package common

func GetValuesFile() string {
	env := GetEnv()
	values := "values.yaml"
	switch env {
	case CI:
		values = "values-ci.yaml"
	case INT:
		values = "values.yaml"
	case PROD:
		values = "values-prod-cus.yaml"
	default:
		values = "values.yaml"
	}
	return values
}

func GetKeyVaultName() string {
	env := GetEnv()
	keyvaultName := "ps2ckvci"
	switch env {
	case CI:
		keyvaultName = "ps2ckvci"
	case INT:
		keyvaultName = "xpay-kv-int"
	case PROD:
		keyvaultName = "xpay-kv-prod"
	default:
		keyvaultName = "ps2ckvci"
	}
	return keyvaultName
}

func GetFlinkHistoryServerURL() string {
	env := GetEnv()
	flinkHistoryServerURL := "https://flink-history.webxtsvc-int.microsoft.com/"
	switch env {
	case CI:
		flinkHistoryServerURL = "https://flink-history.webxtsvc-int.microsoft.com/"
	case INT:
		flinkHistoryServerURL = "https://flink-history.webxtsvc-ppe.microsoft.com/"
	case PROD:
		flinkHistoryServerURL = "https://flink-history-cus.webxtsvc.microsoft.com/"
	default:
		flinkHistoryServerURL = "https://flink-history.webxtsvc-int.microsoft.com/"
	}
	return flinkHistoryServerURL
}

func GetHadoopAzureBlobFileSystemSavepointDirectory() string {
	env := GetEnv()
	hadoopAzureBlobFileSystemSavepointDirectory := ""
	switch env {
	case CI:
		hadoopAzureBlobFileSystemSavepointDirectory = "abfss://xpayflink@xpayadlsgen2ci.dfs.core.windows.net/savepoints/"
	case INT:
		hadoopAzureBlobFileSystemSavepointDirectory = "abfss://xpayflink@xpayadlsgen2int.dfs.core.windows.net/savepoints/"
	case PROD:
		hadoopAzureBlobFileSystemSavepointDirectory = "abfss://xpayflink@xpayadlsgen2prod.dfs.core.windows.net/savepoints/"
	default:
		hadoopAzureBlobFileSystemSavepointDirectory = "abfss://xpayflink@xpayadlsgen2ci.dfs.core.windows.net/savepoints/"
	}
	return hadoopAzureBlobFileSystemSavepointDirectory
}
