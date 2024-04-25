package controller

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"

	"k8s-delegater/pkg/common"
	"k8s-delegater/pkg/common/response"
	"k8s-delegater/pkg/service/html"
)

func CreateTableHTML(ginContext *gin.Context) {

	table := &html.HTMLTable{
		"ADF Job Details",
		[]string{"h1", "h2"},
		[][]string{{"c", "d"}, {"a", "b"}},
	}

	content, err := html.GetTableHTML(table)
	if err != nil {
		common.GetLog().Errorf("html GetTableHTML error: %v", err)
		panic(response.Fail(http.StatusInternalServerError, err.Error(), nil))
	}

	common.GetLog().Errorf("content: %v", content)

	bodyMap := map[string]string{
		"CC":      "xpaydw@microsoft.com,xpayment@microsoft.com",
		"Content": content,
		"Subject": "[ADFExt]ADF Job Details",
		"To":      "miwen@microsoft.com,sallyxiong@microsoft.com,junhua@microsoft.com",
	}
	bodyBytes, err := json.Marshal(bodyMap)
	if err != nil {
		common.GetLog().Errorf("json Marshal error: %v", err)
		panic(response.Fail(http.StatusInternalServerError, err.Error(), nil))
	}

	url, err := html.GetEmailURL()
	if err != nil {
		common.GetLog().Errorf("html GetEmailURL error: %v", err)
		panic(response.Fail(http.StatusInternalServerError, err.Error(), nil))
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		common.GetLog().Errorf("http NewRequest error: %v", err)
		panic(response.Fail(http.StatusInternalServerError, err.Error(), nil))
	}

	req.Header.Set("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		common.GetLog().Errorf("http Call error: %v", err)
		panic(response.Fail(http.StatusInternalServerError, err.Error(), nil))
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		common.GetLog().Errorf("http Call StatusCode: %v", res.StatusCode)
		panic(response.Fail(http.StatusInternalServerError, err.Error(), nil))
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		common.GetLog().Errorf("ioutil ReadAll error: %v", err)
		panic(response.Fail(http.StatusInternalServerError, err.Error(), nil))
	}

	ginContext.JSON(http.StatusOK, response.Success("success", resBody))
}
