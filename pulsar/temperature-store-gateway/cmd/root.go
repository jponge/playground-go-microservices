package cmd

import (
	"github.com/jponge/playground-go-microservices/pulsar/temperature-store-gateway/gateway"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "temperature-store-gateway",
	Short: "Temperature store API that acts as a gateway to Pulsar",
	Run: func(cmd *cobra.Command, args []string) {
		service := gateway.Service{
			Host:        viper.GetString("http.host"),
			Port:        viper.GetInt("http.port"),
			PulsarURL:   viper.GetString("pulsar.broker-address"),
			PulsarTopic: viper.GetString("pulsar.topic"),
		}
		service.Start()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "config.yaml", "config file (default is config.yaml)")

	rootCmd.Flags().String("http.host", "localhost", "Server host")
	rootCmd.Flags().Int("http.port", 3000, "Server port")

	rootCmd.Flags().String("pulsar.broker-address", "pulsar://localhost:6650", "Pulsar broker address")
	rootCmd.Flags().String("pulsar.topic", "temperature-updates", "Target consumer topic")
}

func initConfig() {
	viper.SetConfigFile(cfgFile)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		log.Println("💡 Using config file:", viper.ConfigFileUsed())
	}

	err := viper.BindPFlags(rootCmd.Flags())
	if err != nil {
		log.Fatalln(err)
	}
}
