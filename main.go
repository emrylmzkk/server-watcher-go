package main

import (
	"context"
	"log"
	"os"
	"server-watcher-app/src"
	"server-watcher-app/src/generic"

	"github.com/gofiber/fiber/v2"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()

	if err != nil {
		log.Println(".env variable not uploaded")
	}

	db, err := generic.InitDB()

	if err != nil {
		log.Fatal("Db connection error", err)
	}

	app := fiber.New()
	app.Use(requestid.New())
	app.Use(fiberLogger.New(fiberLogger.Config{
		Format: "[${time}] ${status} ${method} ${path} | ${latency} | id=${locals:requestid} err=${error}\n",
	}))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	container := src.NewAppContainer(db)

	if err := container.InitHelper.CreateAdminUser(ctx); err != nil {
		log.Println("Admin createion failed...", err)
	}

	if err := container.InitHelper.CreateMockPm2Project(ctx); err != nil {
		log.Println("Mock pm2 project creation failed...", err)
	}

	//container.SyncWorker.Start(ctx)
	container.ContainerStatsWorker.Start(ctx)

	container.SystemNotification.Start(ctx)

	src.SetupRoutes(app, container)

	backendPort := os.Getenv("BACKEND_PORT")

	log.Fatal(app.Listen(":" + backendPort))
}
