package routes

import (
	"github.com/gin-gonic/gin"
	"uam-power-backend/controller/aircraft_id_controller"
	"uam-power-backend/models/config_models/db_config_model"
	"uam-power-backend/utils"
)

func SetupAircraftIdRoutes(
	r *gin.Engine, RedisCfg *db_config_model.RedisConfigModel,
	MySqlCfg *db_config_model.MySqlConfigModel,
) {
	aircraftIDController := aircraft_id_controller.NewAircraftIdController(MySqlCfg, RedisCfg)
	uploadApis := r.Group("/aircraftID")
	uploadApis.POST("/info", aircraftIDController.GetAircraftInfo)
	uploadApis.POST("/create", aircraftIDController.CreateUser)
	utils.MsgSuccess("    [AircraftIdRoutes]Successfully init!")
}
