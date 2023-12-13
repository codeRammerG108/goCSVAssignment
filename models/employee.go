package models

import "gorm.io/gorm"

type Employee struct {
	gorm.Model
	ID    string
	Name  string
	Email string
}
