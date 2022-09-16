package main

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jponge/playground-go-microservices/simple-temperature-store/model"
	"github.com/stretchr/testify/assert"
)

func TestAPI(t *testing.T) {
	app := setupFiberApp()

	type args struct {
		method         string
		path           string
		payload        io.Reader
		expectedStatus int
		expectedJSON   interface{}
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
				expectedJSON: []model.TemperatureUpdate{
					{
						SensorID: "123-abc",
						Value:    19.2,
					},
					{
						SensorID: "456-def",
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
				expectedJSON: model.TemperatureUpdate{
					SensorID: "1",
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
				expectedJSON: model.TemperatureUpdate{
					SensorID: "1",
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
				expectedJSON:   nil,
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

			if tt.args.expectedJSON == nil {
				return
			}
			expectedJSON, err := json.Marshal(tt.args.expectedJSON)
			assert.NoError(t, err, "JSON marshalling")
			actualJSON, err := io.ReadAll(resp.Body)
			assert.NoError(t, err, "Body response")
			assert.JSONEq(t, string(expectedJSON), string(actualJSON))
		})
	}
}
