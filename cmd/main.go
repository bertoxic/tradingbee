package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/bertoxic/tradingbee/configs"
	handler "github.com/bertoxic/tradingbee/internal/handlers"
	"github.com/bertoxic/tradingbee/internal/transport/http"
	"github.com/bertoxic/tradingbee/pkg/database"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file",err)
	}
	app, err := config.NewAppConfig()
	if err != nil {
		log.Println("error occured while setting env", err)
	}
    port := os.Getenv("PORT")
    if port == "" {
        port = "9000" // Default to 5001 if PORT is not set
    }
	portnum, err := strconv.Atoi(port)


	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	client := database.NewMongodClient(ctx, app.Config.MONGODB_URI)
	app.InProduction = false
	app.Client = client
	handler.NewDataB(app)
	
	routes := routes.Router()

	if err != nil {
		log.Println("error getting or converting port")
	}
	srv := &http.Server{
		Addr:    ":" + fmt.Sprintf("%d", portnum),
		Handler: routes,
	}
	err = srv.ListenAndServe()
	if err != nil {
        log.Println("vvvvvvvv",portnum)
		log.Println("Eroooxxx " + err.Error())
		return
	}
   

}
