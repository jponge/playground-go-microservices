package cmd

import (
	"github.com/jponge/playground-go-microservices/db-temperature-store/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the server",
	Run: func(cmd *cobra.Command, args []string) {
		dbType := viper.GetString("db")
		if dbType == "sqlite" {
			dbFile := viper.GetString("db.sqlite.file")
			server.InitDb(sqlite.Open(dbFile), &gorm.Config{})
		} else {
			log.Fatal("DB type not supported", dbType)
		}
		server.Start(viper.GetString("http.host"), viper.GetInt("http.port"))
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().String("http.host", "localhost", "Server host")
	runCmd.Flags().Int("http.port", 3000, "Server port")
	runCmd.Flags().String("db", "sqlite", "Database")
	runCmd.Flags().String("db.sqlite.file", "data.db", "Database file (for sqlite)")

	err := viper.BindPFlags(runCmd.Flags())
	if err != nil {
		log.Fatalln(err)
	}
}
