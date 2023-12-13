package main

import (
	db "gocsv/db"
	handlers "gocsv/handlers"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {

	// fmt.Println("The server ")
	app := fiber.New()

	// Initialize database
	db.InitDatabase()

	// Employee handlers
	app.Post("/employees/csv", handlers.CreateEmployeeFromCSV)
	app.Put("/employees/csv", handlers.UpdateEmployeeFromCSV)
	app.Get("/employees/:id/csv", handlers.GetEmployeeByIDCSV)
	app.Get("/employees/csv", handlers.GetAllEmployeesCSV)
	app.Delete("/employees/:id", handlers.DeleteEmployeeByID)

	log.Fatal(app.Listen(":3000"))
}
