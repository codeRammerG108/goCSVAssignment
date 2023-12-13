package handlers

import (
	"encoding/csv"
	"fmt"
	db "gocsv/db"
	models "gocsv/models"
	"io"
	"reflect"

	"github.com/gofiber/fiber/v2"
)

func CreateEmployeeFromCSV(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	csvfile, err := file.Open()
	if err != nil {
		return err
	}
	defer csvfile.Close()

	reader := csv.NewReader(csvfile)
	header, err := reader.Read()
	if err != nil && err != io.EOF {
		return err
	}

	expectedHeader := []string{"ID", "Name", "Email"}
	if !reflect.DeepEqual(header, expectedHeader) {
		return fmt.Errorf("invalid CSV format. Expected header: %v", expectedHeader)
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if len(record) != len(expectedHeader) {
			return fmt.Errorf("invalid CSV record format. Expected %d columns", len(expectedHeader))
		}

		employee := models.Employee{
			ID:    record[0],
			Name:  record[1],
			Email: record[2],
		}
		db.DB.Create(&employee)
	}

	return c.SendString("Employees created successfully from CSV")
}

// func CreateEmployeeFromCSV(c *fiber.Ctx) error {
// 	file, err := c.FormFile("file")
// 	if err != nil {
// 		return err
// 	}

// 	csvfile, err := file.Open()
// 	if err != nil {
// 		return err
// 	}
// 	defer csvfile.Close()

// 	reader := csv.NewReader(csvfile)
// 	// Skip header row
// 	_, err = reader.Read()
// 	if err != nil && err != io.EOF {
// 		return err
// 	}

// 	for {
// 		record, err := reader.Read()
// 		if err == io.EOF {
// 			break
// 		}
// 		if err != nil {
// 			return err
// 		}

// 		employee := models.Employee{
// 			ID:    record[0],
// 			Name:  record[1],
// 			Email: record[2],
// 		}
// 		db.DB.Create(&employee)
// 	}

// 	return c.SendString("Employees created successfully from CSV")
// }

func UpdateEmployeeFromCSV(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	csvfile, err := file.Open()
	if err != nil {
		return err
	}
	defer csvfile.Close()

	reader := csv.NewReader(csvfile)
	// Skip header row
	_, err = reader.Read()
	if err != nil && err != io.EOF {
		return err
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		var employee models.Employee
		if err := db.DB.First(&employee, record[0]).Error; err != nil {
			return fmt.Errorf("employee ID %s not found", record[0])
		}

		employee.Name = record[1]
		employee.Email = record[2]

		db.DB.Save(&employee)
	}

	return c.SendString("Employees updated successfully from CSV")
}
