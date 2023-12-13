package main

import (
	"fmt"
	db "gocsv/db"
	handlers "gocsv/handlers"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {

	fmt.Println("The server is getting initialized")
	app := fiber.New()

	// Initialize database
	db.InitDatabase()

	// Employee handlers
	app.Post("/employees/csv", handlers.CreateEmployeeFromCSV)
	app.Put("/employees/csv", handlers.UpdateEmployeeFromCSV)

	log.Fatal(app.Listen(":3000"))
}
