package main

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/xtt28/shortener/api/router"
	"github.com/xtt28/shortener/database"
)

func main() {
	database.InitDBOrPanic()
	database.MigrateAllModels(database.DB)
	router.InitAndStartRouter()
}
