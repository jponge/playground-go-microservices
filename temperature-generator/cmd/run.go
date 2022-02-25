package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"io"
	"math/rand"
	"net/http"
	"time"
)

func init() {
	runCmd.Flags().String("host", "localhost", "The target temperature store host")
	runCmd.Flags().Uint("port", 3000, "The target temperature store port")
	runCmd.Flags().Uint("count", 5, "The number of generators")
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the generator",
	Run: func(cmd *cobra.Command, args []string) {
		targetHost, err := cmd.Flags().GetString("host")
		if err != nil {
			panic(err)
		}
		targetPort, err := cmd.Flags().GetUint("port")
		if err != nil {
			panic(err)
		}
		count, err := cmd.Flags().GetUint("count")
		if err != nil {
			panic(err)
		}
		httpClient := &http.Client{}
		url := fmt.Sprintf("http://%s:%d/record", targetHost, targetPort)
		fmt.Println("🚀 Start with endpoint", url)
		for i := 0; i < int(count); i++ {
			go runGenerator(url, httpClient)
		}
		runStarted = true
	},
}

func runGenerator(url string, httpClient *http.Client) {
	sensorId := uuid.NewString()
	var temperature = 21.0
	for range time.Tick(5 * time.Second) {
		time.Sleep(time.Duration(rand.Intn(3000)) * time.Millisecond)
		temperature = updateTemperature(temperature)
		fmt.Println("tick", sensorId, temperature)
		jsonData := &payload{SensorId: sensorId, Value: temperature}
		err := performHttpRequest(url, httpClient, jsonData)
		if err != nil {
			fmt.Println("HTTP request failed", err)
		}
	}
}

func performHttpRequest(url string, httpClient *http.Client, data *payload) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	request.Header.Add("Content-Type", "application/json")
	response, err := httpClient.Do(request)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)
	if response.StatusCode != 200 {
		return errors.New(fmt.Sprintf("HTTP request did not result in a 200 status code but %d", response.StatusCode))
	}
	return nil
}

type payload struct {
	SensorId string  `json:"sensorId"`
	Value    float64 `json:"value"`
}

func updateTemperature(temperature float64) float64 {
	delta := rand.Float64() / 10.0
	if rand.Int()%2 == 0 {
		delta = -delta
	}
	temperature = temperature + delta
	return temperature
}
