package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jponge/playground-go-microservices/simple-temperature-store/model"
	"log"
)

func AllDataHandler(db *model.Database) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		log.Printf("Data request from %s", ctx.Request().Host())
		return ctx.JSON(db.AllTemperatureUpdates())
	}
}

func SingleDataHandler(db *model.Database) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		log.Printf("Data request from %s of %s", ctx.Request().Host(), id)
		update, err := db.OneTemperatureUpdate(id)
		if err != nil {
			return ctx.Status(404).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return ctx.JSON(update)
	}
}

func RecordHandler(db *model.Database) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		update := new(model.TemperatureUpdate)
		if err := ctx.BodyParser(&update); err != nil {
			log.Println(err, "->", string(ctx.Body()))
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Could not parse JSON from the request body",
			})
		}
		log.Printf("Record request from %s with payload %s", ctx.Request().Host(), update)
		db.Put(update.SensorID, update.Value)
		return ctx.JSON(update)
	}
}
