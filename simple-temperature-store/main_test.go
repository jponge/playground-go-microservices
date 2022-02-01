package main

import (
	"encoding/json"
	"github.com/jponge/playground-go-microservices/simple-temperature-store/model"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAPI(t *testing.T) {
	app := setupFiberApp()

	type args struct {
		method         string
		path           string
		payload        io.Reader
		expectedStatus int
		expectedJson   interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Get the initial data",
			args: args{
				method:         "GET",
				path:           "/data",
				payload:        nil,
				expectedStatus: 200,
				expectedJson: []model.TemperatureUpdate{
					{
						SensorId: "123-abc",
						Value:    19.2,
					},
					{
						SensorId: "456-def",
						Value:    -2.33,
					},
				},
			},
		},
		{
			name: "Post an update",
			args: args{
				method: "POST",
				path:   "/record",
				payload: strings.NewReader(`{
    				"sensorId": "1",
    				"value": 19.1
				}`),
				expectedStatus: 200,
				expectedJson: model.TemperatureUpdate{
					SensorId: "1",
					Value:    19.1,
				},
			},
		},
		{
			name: "Get the new update data",
			args: args{
				method:         "GET",
				path:           "/data/1",
				payload:        nil,
				expectedStatus: 200,
				expectedJson: model.TemperatureUpdate{
					SensorId: "1",
					Value:    19.1,
				},
			},
		},
		{
			name: "404 for unknown sensor",
			args: args{
				method:         "GET",
				path:           "/data/666",
				payload:        nil,
				expectedStatus: 404,
				expectedJson:   nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.args.method, tt.args.path, tt.args.payload)
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)

			assert.NoError(t, err, "No error")
			assert.Equal(t, tt.args.expectedStatus, resp.StatusCode, "Status code")

			if tt.args.expectedJson == nil {
				return
			}
			expectedJson, err := json.Marshal(tt.args.expectedJson)
			assert.NoError(t, err, "JSON marshalling")
			actualJson, err := ioutil.ReadAll(resp.Body)
			assert.NoError(t, err, "Body response")
			assert.JSONEq(t, string(expectedJson), string(actualJson))
		})
	}
}
