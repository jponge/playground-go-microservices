package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/jponge/playground-go-microservices/db-temperature-store/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"net/http"
)

var db *gorm.DB

func init() {
	dbc, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("DB opening failed", err)
	}
	err = dbc.AutoMigrate(&model.TemperatureUpdate{})
	if err != nil {
		log.Fatal("Migration failed", err)
	}
	db = dbc
}

func send500(writer http.ResponseWriter, err string) {
	writer.WriteHeader(500)
	fmt.Fprintf(writer, err)
}

func Record(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	bytes, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println(err)
		return
	}
	update := &model.TemperatureUpdate{}
	err = json.Unmarshal(bytes, update)
	if err != nil {
		log.Println("JSON decoding failed")
		send500(writer, err.Error())
		return
	}
	entity := &model.TemperatureUpdate{}
	result := db.Where("sensor_id = ?", update.SensorId).First(entity)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			result = db.Create(update)
			if result.Error != nil {
				log.Println("Recording failed", result.Error)
				send500(writer, result.Error.Error())
				return
			}
			entity = update
		} else {
			log.Println("Something went wrong")
			send500(writer, result.Error.Error())
			return
		}
	}
	entity.Value = update.Value
	result = db.Save(entity)
	if result.Error != nil {
		log.Println("Updating record failed", result.Error)
		send500(writer, result.Error.Error())
		return
	}
	bytes, err = json.Marshal(entity)
	if err != nil {
		log.Println("JSON encoding failed")
		send500(writer, err.Error())
		return
	}
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(200)
	fmt.Fprintf(writer, string(bytes))
}

func FetchOne(writer http.ResponseWriter, request *http.Request) {
	sensorId := chi.URLParam(request, "id")
	entity := &model.TemperatureUpdate{}
	result := db.Where("sensor_id = ?", sensorId).First(entity)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			writer.WriteHeader(404)
			return
		}
		send500(writer, result.Error.Error())
		return
	}
	bytes, err := json.Marshal(entity)
	if err != nil {
		log.Println("JSON encoding failed")
		send500(writer, err.Error())
		return
	}
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(200)
	fmt.Fprintf(writer, string(bytes))
}

func FetchAll(writer http.ResponseWriter, request *http.Request) {
	var allEntities []model.TemperatureUpdate
	result := db.Find(&allEntities)
	if result.Error != nil {
		send500(writer, result.Error.Error())
		return
	}
	bytes, err := json.Marshal(allEntities)
	if err != nil {
		log.Println("JSON encoding failed")
		send500(writer, err.Error())
		return
	}
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(200)
	fmt.Fprintf(writer, string(bytes))
}
