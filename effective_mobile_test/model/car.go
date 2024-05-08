package model

type Car struct {
	Base
	RegNum   string `gorm:"not null;unique"`
	Mark     string `gorm:"not null"`
	CarModel string `gorm:"not null"`
	Year     int
	PersonID string `gorm:"not null"`
}
