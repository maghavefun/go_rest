package model

type Person struct {
	Base
	Name       string `gorm:"not null"`
	Surname    string `gorm:"not null"`
	Patronymic string
	Cars       []Car `gorm:"foreignKey:PersonID;"`
}
