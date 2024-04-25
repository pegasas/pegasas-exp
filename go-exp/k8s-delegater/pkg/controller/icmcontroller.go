package controller

import (
	"encoding/json"
	"k8s-delegater/pkg/service/icm"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"k8s-delegater/pkg/common"
	"k8s-delegater/pkg/common/request"
	"k8s-delegater/pkg/common/response"
)

func CreateIncident(ginContext *gin.Context) {

	var request request.CreateIncidentRequest
	err := ginContext.BindJSON(&request)
	if err != nil {
		common.GetLog().Errorf("CreateIncidentRequest bindJson error: %v", err)
		panic(response.Fail(http.StatusInternalServerError, err.Error(), nil))
	}

	CheckCreateIncidentRequest(request)

	incident := icm.NewCreateIncident()
	incident["Title"] = request.Title
	incident["Severity"] = request.Severity
	descriptionEntry := incident["DescriptionEntries"].([]interface{})[0].(map[string]interface{})["DescriptionEntry"].(map[string]interface{})
	descriptionEntry["RenderType"] = "Html"
	descriptionEntry["Text"] = request.Content

	incidentId, status := icm.CreateIncident(incident)

	data := make(map[string]interface{})
	data["incidentId"] = incidentId
	data["status"] = status

	ginContext.JSON(http.StatusOK, response.Success("success", data))
}

func MitigateIncident(ginContext *gin.Context) {
	request := ginContext.Query("incidentId")

	incidentId, err := strconv.Atoi(request)
	if err != nil {
		common.GetLog().Errorf("http NewRequest error: %v", err)
		panic(response.Fail(http.StatusInternalServerError, err.Error(), nil))
	}

	incident, err := json.Marshal(icm.NewMitigateIncident())
	if err != nil {
		common.GetLog().Errorf("http NewRequest error: %v", err)
		panic(response.Fail(http.StatusInternalServerError, err.Error(), nil))
	}
	data := icm.MitigateIncident(incidentId, string(incident))

	ginContext.JSON(http.StatusOK, response.Success("success", data))
}

func CheckCreateIncidentRequest(request request.CreateIncidentRequest) {
	if request.Title == "" {
		common.GetLog().Error("title is null!")
		panic(response.Fail(http.StatusInternalServerError, "title is null", nil))
	}

	if request.Content == "" {
		common.GetLog().Error("content is null!")
		panic(response.Fail(http.StatusInternalServerError, "content is null", nil))
	}

	if request.Severity == "" {
		common.GetLog().Error("severity is null!")
		panic(response.Fail(http.StatusInternalServerError, "severity is null", nil))
	}
}
