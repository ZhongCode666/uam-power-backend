package routes

import (
	"github.com/gin-gonic/gin"
	"uam-power-backend/controller/aircraft_task_controller"
	"uam-power-backend/models/config_models/db_config_model"
	"uam-power-backend/utils"
)

func SetupAircraftTaskRoutes(
	r *gin.Engine, RedisCfg *db_config_model.RedisConfigModel,
	MySqlCfg *db_config_model.MySqlConfigModel,
) {
	aircraftTaskController := aircraft_task_controller.NewAircraftTaskModel(RedisCfg, MySqlCfg)
	uploadApis := r.Group("/aircraftTask")
	uploadApis.POST("/end", aircraftTaskController.EndTask)
	uploadApis.POST("/create", aircraftTaskController.CreateTask)
	uploadApis.POST("/check", aircraftTaskController.CheckTaskInfo)
	utils.MsgSuccess("    [SetupAircraftTaskRoutes]Successfully init!")
}
