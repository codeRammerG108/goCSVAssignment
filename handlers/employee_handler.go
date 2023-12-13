package handlers

import (
	"encoding/csv"
	"errors"
	"fmt"
	db "gocsv/db"
	models "gocsv/models"
	"io"
	"regexp"
	"strings"

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

	// Validate header fields
	if !isValidHeader(header) {
		return errors.New("CSV should contain ID, Name, and Email columns")
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if err := validateRecord(record); err != nil {
			return err
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

	header, err := reader.Read()
	if err != nil && err != io.EOF {
		return err
	}
	// Validate header fields
	if !isValidHeader(header) {
		return errors.New("CSV should contain ID, Name, and Email columns")
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if err := validateRecord(record); err != nil {
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
func GetEmployeeByIDCSV(c *fiber.Ctx) error {
	employeeID := c.Params("id")

	var employee models.Employee
	if err := db.DB.First(&employee, employeeID).Error; err != nil {
		return fmt.Errorf("employee ID %s not found", employeeID)
	}

	csvData := [][]string{
		{"ID", "Name", "Email"},
		{employee.ID, employee.Name, employee.Email},
	}

	csvString := convertToCSV(csvData)

	return c.SendString(csvString)
}

func convertToCSV(data [][]string) string {
	var csvString string
	for _, record := range data {
		csvString += fmt.Sprintf("%s\n", strings.Join(record, ","))
	}
	return csvString
}
func GetAllEmployeesCSV(c *fiber.Ctx) error {
	var employees []models.Employee
	if err := db.DB.Find(&employees).Error; err != nil {
		return err
	}

	csvData := [][]string{{"ID", "Name", "Email"}}
	for _, employee := range employees {
		csvData = append(csvData, []string{employee.ID, employee.Name, employee.Email})
	}

	csvString := convertToCSV(csvData)

	return c.SendString(csvString)
}
func DeleteEmployeeByID(c *fiber.Ctx) error {
	employeeID := c.Params("id")

	var employee models.Employee
	if err := db.DB.First(&employee, employeeID).Error; err != nil {
		return fmt.Errorf("employee ID %s not found", employeeID)
	}

	db.DB.Delete(&employee, employeeID)

	return c.SendString("Employee deleted successfully")
}
func isValidHeader(header []string) bool {
	requiredColumns := []string{"ID", "Name", "Email"}
	if len(header) != len(requiredColumns) {
		return false
	}
	for i, col := range header {
		if col != requiredColumns[i] {
			return false
		}
	}
	return true
}

func validateRecord(record []string) error {
	if len(record) != 3 {
		return errors.New("each record should contain ID, Name, and Email")
	}
	if record[1] == "" || record[2] == "" || record[0] == "" {
		return errors.New("name and Email fields cannot be empty")
	}
	if !isValidEmail(record[2]) {
		return errors.New("invalid Email format")
	}
	return nil
}

func isValidEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	return regexp.MustCompile(emailRegex).MatchString(email)
}
