package service

import (
	"strings"

	"github.com/maghavefun/effective_mobile_test/external"
	"github.com/maghavefun/effective_mobile_test/initializer"
	"github.com/maghavefun/effective_mobile_test/model"
	"gorm.io/gorm"
)

type CarsFilter struct {
	RegNum  string
	Mark    string
	Model   string
	Year    string
	OwnerID string
	Page    int
	PerPage int
}

type ICarService interface {
	GetCarsWithFilter(filter CarsFilter) ([]model.Car, error)
	DeleteCarByID(carId string) error
	UpdateCarByID(carId string, carDTO model.Car) (model.Car, error)
	CreateCar(amountOfCars int, carsCh chan external.CarDTO, errorCh chan error) ([]model.Car, error)
}

type CarService struct{}

func NewCarService() ICarService {
	return &CarService{}
}

func (s *CarService) GetCarsWithFilter(filter CarsFilter) ([]model.Car, error) {
	var cars []model.Car

	offset := (filter.Page - 1) * filter.PerPage
	queryDB := initializer.DB.Where("true").Session(&gorm.Session{})
	if offset != 0 && filter.PerPage != 0 {
		queryDB = queryDB.Offset(offset).Limit(filter.PerPage)
	}
	if filter.RegNum != "" {
		queryDB = queryDB.Where("reg_num IN ?", strings.Split(filter.RegNum, ","))
	}
	if filter.Mark != "" {
		queryDB = queryDB.Where("mark ILIKE ?", filter.Mark)
	}
	if filter.Model != "" {
		queryDB = queryDB.Where("car_model ILIKE ?", filter.Model)
	}
	if filter.Year != "" {
		queryDB = queryDB.Where("year = ?", filter.Year)
	}
	if filter.OwnerID != "" {
		queryDB = queryDB.Where("person_id = ?", filter.OwnerID)
	}
	if err := queryDB.Error; err != nil {
		return cars, err
	}
	queryDB.Order("mark DESC").Find(&cars)

	return cars, nil
}

func (s *CarService) DeleteCarByID(carId string) error {
	if err := initializer.DB.Where("id = ?", carId).Delete(&model.Car{}).Error; err != nil {
		return err
	}
	return nil
}

func (s *CarService) UpdateCarByID(carId string, carDTO model.Car) (model.Car, error) {
	var car model.Car
	if err := initializer.DB.Where("id = ?", carId).First(&car).Error; err != nil {
		return car, err
	}

	result := initializer.DB.Model(&car).Where("id = ?", carId).Updates(model.Car{
		RegNum:   carDTO.RegNum,
		Mark:     carDTO.Mark,
		CarModel: carDTO.CarModel,
		Year:     carDTO.Year,
	})

	if result.Error != nil {
		return car, result.Error
	}

	return car, nil
}

func (s *CarService) CreateCar(amountOfCars int, carsCh chan external.CarDTO, errorCh chan error) ([]model.Car, error) {
	var owner model.Person
	var car model.Car
	var cars []model.Car

	for i := 0; i < amountOfCars; i++ {
		select {
		case fetchedCar := <-carsCh:
			err := initializer.DB.Transaction(func(tx *gorm.DB) error {
				owner = model.Person{
					Name:       fetchedCar.Owner.Name,
					Surname:    fetchedCar.Owner.Surname,
					Patronymic: fetchedCar.Owner.Patronymic,
				}
				if err := tx.Create(&owner).Error; err != nil {
					return err
				}
				car = model.Car{
					RegNum:   fetchedCar.RegNum,
					Mark:     fetchedCar.Mark,
					CarModel: fetchedCar.Model,
					Year:     fetchedCar.Year,
					PersonID: owner.ID,
				}

				if err := tx.Create(&car).Error; err != nil {
					return err
				}
				cars = append(cars, car)
				return nil
			})
			if err != nil {
				return cars, err
			}
		case err := <-errorCh:
			return cars, err
		}
	}

	return cars, nil
}
