package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/jponge/playground-go-microservices/simple-temperature-store/handlers"
	"github.com/jponge/playground-go-microservices/simple-temperature-store/model"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

func main() {
	mainCommand := setupCobraAndViper()
	if err := mainCommand.Execute(); err != nil {
		log.Fatal(err)
	}
}

func setupCobraAndViper() cobra.Command {
	mainCommand := cobra.Command{
		Use:     "simple-temperature-store",
		Short:   "A HTTP service to store temperature update data",
		Version: "0.1",
		Run:     start,
	}

	mainCommand.Flags().String("http.host", "localhost", "Host to run the HTTP server")
	mainCommand.Flags().Int("http.port", 3000, "Port to run the HTTP server on")
	if err := viper.BindPFlags(mainCommand.Flags()); err != nil {
		log.Fatal(err)
	}

	viper.SetDefault("http.host", "localhost")
	viper.SetDefault("http.port", 3000)

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.simple-temperature-store")
	viper.AddConfigPath("/etc/simple-temperature-store")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("No config.yaml file was found")
		} else {
			log.Fatal(err)
		}
	}

	if err := viper.BindEnv("http.host", "HTTP_HOST"); err != nil {
		log.Fatal(err)
	}
	if err := viper.BindEnv("http.port", "HTTP_PORT"); err != nil {
		log.Fatal(err)
	}
	return mainCommand
}

func start(cmd *cobra.Command, args []string) {
	app := setupFiberApp()

	host := fmt.Sprintf("%s:%d", viper.Get("http.host"), viper.GetInt("http.port"))
	log.Printf("Listening on http://%s", host)
	err := app.Listen(host)
	if err != nil {
		log.Fatal(err)
	}
}

func setupFiberApp() *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:               "simple-temperature-store",
		DisableStartupMessage: true,
	})

	db := model.NewDatabase()
	db.Put("123-abc", 19.2)
	db.Put("456-def", -2.33)

	app.Get("/data", handlers.AllDataHandler(db))
	app.Get("/data/:id", handlers.SingleDataHandler(db))
	app.Post("/record", handlers.RecordHandler(db))

	return app
}
