package main

import (
	"context"
	"embed"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func init() {
	viper.SetDefault("listen.address", "0.0.0.0")
	viper.SetDefault("listen.port", 4000)
	viper.SetDefault("listen.ssl.certFile", "cert.pem")
	viper.SetDefault("listen.ssl.keyFile", "key.pem")
	viper.SetDefault("api.host", "localhost")
	viper.SetDefault("api.port", 3000)

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.dashboard-webapp")
	viper.AddConfigPath("/etc/dashboard-webapp")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("No config.yaml file was found")
			os.Exit(1)
		} else {
			log.Fatal(err)
		}
	}

	address = fmt.Sprintf("%s:%d", viper.GetString("listen.address"), viper.GetInt("listen.port"))
	apiURL = fmt.Sprintf("http://%s:%d/data", viper.GetString("api.host"), viper.GetInt("api.port"))
}

var address string
var apiURL string

//go:embed assets
var assets embed.FS

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/data", fetchData).Methods("GET")
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/", http.FileServer(http.FS(assets))))
	router.HandleFunc("/", serveRoot).Methods("GET")

	server := &http.Server{
		Addr:         address,
		Handler:      handlers.LoggingHandler(os.Stdout, router),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
	go func() {
		log.Println("🚀 Start to listen on", address)
		err := server.ListenAndServeTLS(viper.GetString("listen.ssl.certFile"), viper.GetString("listen.ssl.keyFile"))
		if err != nil {
			log.Println(err)
		}
	}()

	// Wait for signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	server.Shutdown(ctx)
	log.Println("👋 Bye!")
}

func serveRoot(writer http.ResponseWriter, request *http.Request) {
	http.Redirect(writer, request, "/assets/index.html", http.StatusFound)
}

func fetchData(writer http.ResponseWriter, request *http.Request) {
	apiResponse, err := http.Get(apiURL)
	if err != nil {
		sendHTTP500(writer, err)
		return
	}
	defer apiResponse.Body.Close()
	writer.Header().Add("Content-Type", "application/json")
	data, err := ioutil.ReadAll(apiResponse.Body)
	if err != nil {
		sendHTTP500(writer, err)
	}
	_, err = writer.Write(data)
	if err != nil {
		log.Println(err)
	}
}

func sendHTTP500(writer http.ResponseWriter, err error) {
	log.Println(err)
	writer.WriteHeader(500)
	writer.Write([]byte(err.Error()))
}
