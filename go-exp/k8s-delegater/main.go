package main

import (
	"github.com/gin-gonic/gin"

	"k8s-delegater/pkg/common"
	"k8s-delegater/pkg/controller"
)

func main() {
	common.InitConfig()
	common.InitLogger()

	engine := gin.Default()

	engine.Use(common.Recover)
	engine.Use(common.Auth)

	engine.GET("/ping", controller.Ping)
	engine.GET("/api/keepalive", controller.KeepAlive)
	engine.GET("/probe/healthcheck", controller.ProbeHealthCheck)

	engine.POST("/api/flink/deployment", controller.CreateFlinkDeploymentJob)
	engine.GET("/api/flink/deployment/status", controller.GetFlinkDeploymentStatus)
	engine.GET("/api/flink/deployment/log", controller.GetFlinkDeploymentLog)
	engine.GET("/api/flink/job/id", controller.GetFlinkJobID)
	engine.POST("/api/flink/job/status", controller.GetFlinkJobStatus)
	engine.DELETE("/api/flink/job", controller.DeleteFlinkJob)
	engine.POST("/api/flink/savepoint", controller.TriggerFlinkJobSavepoint)
	engine.POST("/api/flink/savepoint/trigger", controller.GetFlinkJobSavepointTriggerStatus)

	engine.GET("/html", controller.CreateTableHTML)

	engine.POST("/api/incident/create", controller.CreateIncident)
	engine.PUT("/api/incident/mitigate", controller.MitigateIncident)

	engine.Run(":8080")
}
