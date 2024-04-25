package controller

import (
	"github.com/gin-gonic/gin"
	"k8s-delegater/pkg/common/response"
	"net/http"
)

func Ping(ginContext *gin.Context) {
	ginContext.JSON(http.StatusOK, response.Success("Pong", nil))
}

func KeepAlive(ginContext *gin.Context) {
	ginContext.JSON(http.StatusOK, response.Success("", nil))
}

func ProbeHealthCheck(ginContext *gin.Context) {
	ginContext.JSON(http.StatusOK, response.Success("", nil))
}
