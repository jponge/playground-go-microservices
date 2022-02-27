package server

import (
	"github.com/jponge/playground-go-microservices/db-temperature-store/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func setupGormWithSqliteDb() (string, error) {
	file, err := os.CreateTemp("", "test-db")
	if err != nil {
		return "", err
	}
	log.Println("Using temporary DB file", file.Name())
	InitDb(sqlite.Open(file.Name()), &gorm.Config{})
	return file.Name(), nil
}

func TestApiWithSqliteDb(t *testing.T) {
	fileToCleanup, err := setupGormWithSqliteDb()
	if err != nil {
		t.Fatal(err)
		return
	}
	defer os.Remove(fileToCleanup)

	server := httptest.NewServer(AppRouter())
	defer server.Close()

	createUpdate(t, server, &model.TemperatureUpdate{
		SensorId: "123-abc",
		Value:    19.0,
	})
	fetchSensorData(t, server, "123-abc", true, 19.0)
	fetchSensorData(t, server, "foo-bar", false, 0)
	createUpdate(t, server, &model.TemperatureUpdate{
		SensorId: "123-abc",
		Value:    19.2,
	})
	fetchSensorData(t, server, "123-abc", true, 19.2)
}

func createUpdate(t *testing.T, server *httptest.Server, update *model.TemperatureUpdate) {
	jsonReader, err := update.ToJSONReader()
	if err != nil {
		t.Fatal(err)
	}
	response, err := http.Post(server.URL+"/record", "application/json", jsonReader)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	assert.Equal(t, 200, response.StatusCode)
	responseUpdate, err := model.TemperatureUpdateFromJSONReader(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, update.SensorId, responseUpdate.SensorId)
	assert.Equal(t, update.Value, responseUpdate.Value)
}

func fetchSensorData(t *testing.T, server *httptest.Server, sensorId string, expectedToBeFound bool, expectedTemperature float64) {
	response, err := http.Get(server.URL + "/data/" + sensorId)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	expectedStatusCode := 200
	if !expectedToBeFound {
		expectedStatusCode = 404
	}
	assert.Equal(t, expectedStatusCode, response.StatusCode)
	if expectedToBeFound {
		sensorData, err := model.TemperatureUpdateFromJSONReader(response.Body)
		if err != nil {
			t.Fatal(err)
			return
		}
		assert.Equal(t, sensorId, sensorData.SensorId)
		assert.Equal(t, expectedTemperature, sensorData.Value)
	}
}
