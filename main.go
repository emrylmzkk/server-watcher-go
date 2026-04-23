package main

import (
	"log"
	"server-watcher-app/src"
	"server-watcher-app/src/generic"

	"github.com/gofiber/fiber/v2"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func main() {

	db, err := generic.InitDB()

	if err != nil {
		log.Fatal("Db connection error", err)
	}

	app := fiber.New()
	app.Use(requestid.New())
	app.Use(fiberLogger.New(fiberLogger.Config{
		Format: "[${time}] ${status} ${method} ${path} | ${latency} | id=${locals:requestid} err=${error}\n",
	}))

	// routes.SetupAuthRoutes(app, container.AuthHandler)
	// routes.SetupUserRoutes(app, container.UserHandler, container.KeycloakService, &container.Userservice)
	// routes.SetupFamilyRoutes(app, container.FamilyHandler, container.KeycloakService, container.Userservice)

	container := src.NewAppContainer(db)

	src.SetupRoutes(app, container)

	log.Fatal(app.Listen(":8005"))
}
