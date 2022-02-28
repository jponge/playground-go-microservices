package cmd

import (
	"fmt"
	"github.com/jponge/playground-go-microservices/db-temperature-store/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the server",
	Run: func(cmd *cobra.Command, args []string) {
		dbType := viper.GetString("db")
		switch dbType {
		case "sqlite":
			dbFile := viper.GetString("db.sqlite.file")
			server.InitDb(sqlite.Open(dbFile), &gorm.Config{})
		case "postgres":
			dsn := fmt.Sprintf(
				"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Europe/Paris",
				viper.GetString("db.postgres.host"),
				viper.GetString("db.postgres.user"),
				viper.GetString("db.postgres.password"),
				viper.GetString("db.postgres.database"),
				viper.GetInt("db.postgres.port"),
			)
			server.InitDb(postgres.Open(dsn), &gorm.Config{})
		default:
			log.Fatal("DB type not supported", dbType)
		}
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
	runCmd.Flags().String("db", "sqlite", "Database [sqlite, postgres]")
	runCmd.Flags().String("db.sqlite.file", "data.db", "Database file (for sqlite)")
	runCmd.Flags().String("db.postgres.host", "localhost", "Postgres host")
	runCmd.Flags().Int("db.postgres.port", 5432, "Postgres port")
	runCmd.Flags().String("db.postgres.user", "postgres", "Postgres user")
	runCmd.Flags().String("db.postgres.password", "postgres", "Postgres password")
	runCmd.Flags().String("db.postgres.database", "postgres", "Postgres database")

	err := viper.BindPFlags(runCmd.Flags())
	if err != nil {
		log.Fatalln(err)
	}
}
