package aircraft_data_controller

import (
	"github.com/gofiber/fiber/v2"
	"strconv"
	"uam-power-backend/models/config_models/db_config_model"
	"uam-power-backend/models/controller_models/data_flow_model"
	"uam-power-backend/service/db_service"
	"uam-power-backend/utils"
)

type RequestAircraft struct {
	StatusRedisService *dbservice.RedisDict
	EventRedisService  *dbservice.RedisDict
}

func NewReceiveAircraft(redisConfig *db_config_model.RedisConfigModel) *RequestAircraft {
	utils.MsgSuccess("        [ReceiveAircraft]init successfully!")
	redisStatusService := dbservice.NewRedisDict(redisConfig.Host, redisConfig.Port, redisConfig.StatusDBno)
	redisEventService := dbservice.NewRedisDict(redisConfig.Host, redisConfig.Port, redisConfig.EventDBno)
	return &RequestAircraft{StatusRedisService: redisStatusService, EventRedisService: redisEventService}
}

func (receiver *RequestAircraft) RequestAircraftStatus(c *fiber.Ctx) error {
	var aircraftReq data_flow_model.RecAircraftStatusRequest

	// 绑定 JSON 数据到结构体
	if err := c.BodyParser(&aircraftReq); err != nil {
		utils.MsgError("        [ReceiveAircraft]RequestAircraftStatus Invalid JSON data! >" + err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid JSON data"})
	}

	rec, err := receiver.StatusRedisService.Get(strconv.Itoa(aircraftReq.AircraftID))
	if err != nil || rec == nil {
		utils.MsgError("        [ReceiveAircraft]RequestAircraftStatus Invalid JSON data!")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"msg": "N.A.!"})
	}
	utils.MsgError("        [ReceiveAircraft]RequestAircraftStatus Invalid JSON data!")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"msg": "Successfully requestData!", "data": rec})
}

func (receiver *RequestAircraft) RequestAircraftEvent(c *fiber.Ctx) error {
	var aircraftReq data_flow_model.RecAircraftStatusRequest

	// 绑定 JSON 数据到结构体
	if err := c.BodyParser(&aircraftReq); err != nil {
		utils.MsgError("        [ReceiveAircraft]RequestAircraftEvent Invalid JSON data! >" + err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid JSON data"})
	}

	rec, err := receiver.EventRedisService.Get(strconv.Itoa(aircraftReq.AircraftID))
	if err != nil || rec == nil {
		utils.MsgError("        [ReceiveAircraft]RequestAircraftEvent Invalid JSON data!")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"msg": "N.A.!"})
	}
	utils.MsgSuccess("        [ReceiveAircraft]RequestAircraftEvent Successfully requestData!")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"msg": "Successfully requestData!", "data": rec})
}
