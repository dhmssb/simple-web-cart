package app

import (
	"flag"
	"log"
	"os"
	"simpleWebCart/app/controllers"

	"github.com/joho/godotenv"
)

func Run() {
	var (
		server = controllers.Server{}
		dbConf = controllers.DBConfig{}
	)

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error load .env file")
	}

	flag.Parse()
	arg := flag.Arg(0)
	if arg != "" {
		server.InitCommands(dbConf)
	} else {
		server.Initialize(dbConf)
		server.Run(os.Getenv("APP_PORT"))
	}
}
