package cmd

import (
	"github.com/jponge/playground-go-microservices/db-temperature-store/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the server",
	Run: func(cmd *cobra.Command, args []string) {
		server.Start(viper.GetString("http.host"), viper.GetInt("http.port"))
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().String("http.host", "localhost", "Server host")
	runCmd.Flags().Int("http.port", 3000, "Server port")

	err := viper.BindPFlags(runCmd.Flags())
	if err != nil {
		log.Fatalln(err)
	}
}
