package routes

import (
	"github.com/gin-gonic/gin"
	"uam-power-backend/controller/data_controller"
	"uam-power-backend/models/config_models/db_config_model"
	"uam-power-backend/utils"
)

// SetupDataFlowRoutes 配置所有路由
func SetupDataFlowRoutes(
	r *gin.Engine, kafkaCfg *db_config_model.KafkaConfigModel,
	redisCfg *db_config_model.RedisConfigModel,
) {
	aircraftUploadController := data_controller.NewUploadAircraftController(kafkaCfg)
	aircraftReqController := data_controller.NewReceiveAircraft(redisCfg)
	// 设置公共路由
	r.GET("/alive", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	})

	uploadApis := r.Group("/upload")
	uploadApis.POST("/aircraftData", aircraftUploadController.UploadData)
	uploadApis.POST("/aircraftEvent", aircraftUploadController.UploadEvent)

	recApis := r.Group("/request")
	recApis.POST("/aircraftData", aircraftReqController.RequestAircraftStatus)
	recApis.POST("/aircraftEvent", aircraftReqController.RequestAircraftEvent)
	utils.MsgSuccess("    [SetupDataFlowRoutes]Successfully init!")
}
