package server

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/jponge/playground-go-microservices/db-temperature-store/model"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type Controller struct {
	DB *gorm.DB
}

func NewController(dialector gorm.Dialector, config *gorm.Config) Controller {
	dbc, err := gorm.Open(dialector, config)
	if err != nil {
		log.Fatal("DB opening failed: ", err)
	}
	err = dbc.AutoMigrate(&model.TemperatureUpdate{})
	if err != nil {
		log.Fatal("Migration failed: ", err)
	}
	return Controller{
		DB: dbc,
	}
}

func send500(writer http.ResponseWriter, err string) {
	writer.WriteHeader(500)
	writer.Write([]byte(err))
}

func (controller *Controller) Record(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	updateRequest, err := model.TemperatureUpdateFromJSONReader(request.Body)
	if err != nil {
		log.Println("JSON decoding failed")
		send500(writer, err.Error())
		return
	}
	entity := &model.TemperatureUpdate{}
	result := controller.DB.Where("sensor_id = ?", updateRequest.SensorID).First(entity)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			result = controller.DB.Create(updateRequest)
			if result.Error != nil {
				log.Println("Recording failed", result.Error)
				send500(writer, result.Error.Error())
				return
			}
			entity = updateRequest
		} else {
			log.Println("Something went wrong")
			send500(writer, result.Error.Error())
			return
		}
	}
	entity.Value = updateRequest.Value
	result = controller.DB.Save(entity)
	if result.Error != nil {
		log.Println("Updating record failed", result.Error)
		send500(writer, result.Error.Error())
		return
	}
	responseBytes, err := entity.ToJSON()
	if err != nil {
		log.Println("JSON encoding failed")
		send500(writer, err.Error())
		return
	}
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(200)
	writer.Write(responseBytes)
}

func (controller *Controller) FetchOne(writer http.ResponseWriter, request *http.Request) {
	sensorID := chi.URLParam(request, "id")
	entity := &model.TemperatureUpdate{}
	result := controller.DB.Where("sensor_id = ?", sensorID).First(entity)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			writer.WriteHeader(404)
			return
		}
		send500(writer, result.Error.Error())
		return
	}
	responseBytes, err := entity.ToJSON()
	if err != nil {
		log.Println("JSON encoding failed")
		send500(writer, err.Error())
		return
	}
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(200)
	writer.Write(responseBytes)
}

func (controller *Controller) FetchAll(writer http.ResponseWriter, request *http.Request) {
	var allEntities []model.TemperatureUpdate
	result := controller.DB.Find(&allEntities)
	if result.Error != nil {
		send500(writer, result.Error.Error())
		return
	}
	responseBytes, err := json.Marshal(allEntities)
	if err != nil {
		log.Println("JSON encoding failed")
		send500(writer, err.Error())
		return
	}
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(200)
	writer.Write(responseBytes)
}
