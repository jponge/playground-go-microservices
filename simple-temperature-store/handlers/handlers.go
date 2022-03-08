package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jponge/playground-go-microservices/simple-temperature-store/model"
	"log"
)

type API struct {
	DB *model.Database
}

func (api *API) AllDataHandler(ctx *fiber.Ctx) error {
	log.Printf("Data request from %s", ctx.Request().Host())
	return ctx.JSON(api.DB.AllTemperatureUpdates())
}

func (api *API) SingleDataHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	log.Printf("Data request from %s of %s", ctx.Request().Host(), id)
	update, err := api.DB.OneTemperatureUpdate(id)
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.JSON(update)
}

func (api *API) RecordHandler(ctx *fiber.Ctx) error {
	update := new(model.TemperatureUpdate)
	if err := ctx.BodyParser(&update); err != nil {
		log.Println(err, "->", string(ctx.Body()))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Could not parse JSON from the request body",
		})
	}
	log.Printf("Record request from %s with payload %s", ctx.Request().Host(), update)
	api.DB.Put(update.SensorID, update.Value)
	return ctx.JSON(update)
}
