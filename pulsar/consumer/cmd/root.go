package cmd

import (
	"github.com/jponge/playground-go-microservices/pulsar/consumer/consumer"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "consumer",
	Short: "Consume temperature updates from Pulsar",
	Run: func(cmd *cobra.Command, args []string) {
		consumer := consumer.Consumer{
			PulsarURL:   viper.GetString("pulsar.broker-address"),
			PulsarTopic: viper.GetString("pulsar.topic"),
		}
		consumer.Start()
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
