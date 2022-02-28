package server

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/jponge/playground-go-microservices/db-temperature-store/model"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// ------------------------------------------------------------------------- //

func TestAPIWithSqliteDB(t *testing.T) {
	fileToCleanup, err := os.CreateTemp("", "test-db")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(fileToCleanup.Name())

	log.Println("Using temporary DB file", fileToCleanup.Name())
	InitDb(sqlite.Open(fileToCleanup.Name()), &gorm.Config{})

	server := httptest.NewServer(AppRouter())
	defer server.Close()

	performInteractions(t, server)
}

// ------------------------------------------------------------------------- //

func TestAPIWithPostgres(t *testing.T) {
	ctx := context.Background()
	pgPort, err := nat.NewPort("tcp", "5432")
	req := testcontainers.ContainerRequest{
		Image:        "docker.io/postgres:14",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "gorm",
			"POSTGRES_PASSWORD": "gorm",
			"POSTGRES_DB":       "gorm",
		},
		WaitingFor: wait.ForListeningPort(pgPort),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer container.Terminate(ctx)

	mappedHost, err := container.Host(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if err != nil {
		t.Fatal(err)
	}
	mappedPort, err := container.MappedPort(ctx, pgPort)
	if err != nil {
		t.Fatal(err)
	}

	dsn := fmt.Sprintf("host=%s user=gorm password=gorm dbname=gorm port=%d sslmode=disable TimeZone=Europe/Paris", mappedHost, mappedPort.Int())
	log.Println("Connection parameters", dsn)
	InitDb(postgres.Open(dsn), &gorm.Config{})

	server := httptest.NewServer(AppRouter())
	defer server.Close()

	performInteractions(t, server)
}

// ------------------------------------------------------------------------- //

func performInteractions(t *testing.T, server *httptest.Server) {
	log.Println("---- Step 1 ----")

	createUpdate(t, server, &model.TemperatureUpdate{
		SensorID: "123-abc",
		Value:    19.0,
	})

	log.Println("---- Step 2 ----")

	fetchSensorData(t, server, "123-abc", true, 19.0)

	log.Println("---- Step 3 ----")

	fetchSensorData(t, server, "foo-bar", false, 0)

	log.Println("---- Step 4 ----")

	createUpdate(t, server, &model.TemperatureUpdate{
		SensorID: "123-abc",
		Value:    19.2,
	})

	log.Println("---- Step 5 ----")

	fetchSensorData(t, server, "123-abc", true, 19.2)

	log.Println("----  Done  ----")
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
	assert.Equal(t, update.SensorID, responseUpdate.SensorID)
	assert.Equal(t, update.Value, responseUpdate.Value)
}

func fetchSensorData(t *testing.T, server *httptest.Server, sensorID string, expectedToBeFound bool, expectedTemperature float64) {
	response, err := http.Get(server.URL + "/data/" + sensorID)
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
		assert.Equal(t, sensorID, sensorData.SensorID)
		assert.Equal(t, expectedTemperature, sensorData.Value)
	}
}

// ------------------------------------------------------------------------- //
