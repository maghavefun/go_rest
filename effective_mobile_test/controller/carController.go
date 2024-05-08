package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/maghavefun/effective_mobile_test/external"
	"github.com/maghavefun/effective_mobile_test/model"
	"github.com/maghavefun/effective_mobile_test/service"

	"gorm.io/gorm"
)

type carsController struct {
	carService service.ICarService
}

type ICarsController interface {
	GetCars(c *gin.Context)
	Delete(c *gin.Context)
	Update(c *gin.Context)
	Create(c *gin.Context)
}

func NewCarController() *carsController {
	return &carsController{carService: service.NewCarService()}
}

type CreateCarBody struct {
	RegNums []string `json:"regNums"`
}

type UpdateCarBody struct {
	RegNum   string `json:"regnum" example:"VO555X"`
	Mark     string `json:"mark" example:"Toyota"`
	CarModel string `json:"model" example:"Supra"`
	Year     int    `json:"year"  example:"1998"`
	OwnerID  string `json:"ownerId" example:"225a7660-8dff-4d22-93f2-a50606b8ebe6"`
}

// Recieving cars
// @Summary 		Recive cars
// @Description Recive cars list with provided filters and pagination
// @Tags				Cars
// @Produce 		json
// @Param 			perPage   query		 int				false	"page size"
// @Param 			page      query    int				false	"current page"
// @Param 			regNum		query 	 string			false "car plate number"
// @Param 			mark			query		 string			false	"mark of car"
// @Param 			model			query		 string			false	"model of car"
// @Param 			year			query		 int				false "produced year of car"
// @Param 			ownerId		query		 string			false	"id of person who own the car"
// @Success			200				{array}  model.Car
// @Failure     404       {object} map[string]any
// @Failure     500       {object} map[string]any
// @Router      /cars [get]
func (con *carsController) GetCars(c *gin.Context) {
	perPageString := c.Request.URL.Query().Get("perPage")
	perPage, err := strconv.Atoi(perPageString)
	if err != nil && perPageString != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "perPage property should be integer as string",
		})
		return
	}
	pageString := c.Request.URL.Query().Get("page")
	page, err := strconv.Atoi(pageString)
	if err != nil && pageString != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "page property should be integer as string",
		})
		return
	}
	regNumString := c.Request.URL.Query().Get("regNum")
	mark := c.Request.URL.Query().Get("mark")
	carModel := c.Request.URL.Query().Get("model")
	year := c.Request.URL.Query().Get("year")
	ownerId := c.Request.URL.Query().Get("ownerId")

	carsFilter := service.CarsFilter{
		RegNum:  regNumString,
		Mark:    mark,
		Model:   carModel,
		Year:    year,
		OwnerID: ownerId,
		Page:    page,
		PerPage: perPage,
	}

	var cars []model.Car
	cars, err = con.carService.GetCarsWithFilter(carsFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error occured on getting cars",
			"error":   err,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"cars": cars,
	})
}

// Delete car
// @Summary     Delete car
// @Description Delete car by id from param
// @Tags        Cars
// @Param 			id				 path		  string  						true	"car id"
// @Success     200				 {object}	map[string]any
// @Failure     404        {object} map[string]any
// @Router      /cars/{id} [delete]
func (con *carsController) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := con.carService.DeleteCarByID(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("car with id %s not found", id),
			"error":   err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("deleted car with id: %s", id),
	})
}

// Update car
// @Summary     Update car
// @Description Update car by id from param
// @Tags        Cars
// @Accept			json
// @Produce			json
// @Param 			id				path		 string    					true    "car id"
// @Param				car				body		 UpdateCarBody			true    "fields for updating car"
// @Success     200				{object} model.Car
// @Failure     404       {object} map[string]any
// @Failure     500       {object} map[string]any
// @Router      /cars/{id} [put]
func (con *carsController) Update(c *gin.Context) {
	carId := c.Param("id")
	var body UpdateCarBody

	c.Bind(&body)
	var car model.Car
	carDTO := model.Car{
		RegNum:   body.RegNum,
		Mark:     body.Mark,
		CarModel: body.CarModel,
		Year:     body.Year,
		PersonID: body.OwnerID,
	}

	car, err := con.carService.UpdateCarByID(carId, carDTO)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": fmt.Sprintf("Car with id %s not found", carId),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("Cannot delete car with id %s", carId),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("succesfully updated car with id: %s", carId),
		"car":     car,
	})
}

// Create car
// @Summary     Create car
// @Description Create car by regNums param, that provided in body
// @Tags        Cars
// @Accept			json
// @Produce			json
// @Param				regNums		body		 CreateCarBody						true    "body that contains array of plate numbers of car"
// @Success     200				{object} model.Car
// @Failure     404       {object} map[string]any
// @Failure     500       {object} map[string]any
// @Router      /cars [post]
func (con *carsController) Create(c *gin.Context) {
	var body CreateCarBody

	c.Bind(&body)

	carsExternalApi := external.NewCarsAPI()

	carsCh := make(chan external.CarDTO)
	errorCh := make(chan error)
	amountOfCarsToCreate := len(body.RegNums)
	for _, regNum := range body.RegNums {
		queryString := "regNum=" + regNum
		go external.FetchData(carsExternalApi, "cars", queryString, carsCh, errorCh)
	}

	var cars []model.Car

	cars, err := con.carService.CreateCar(amountOfCarsToCreate, carsCh, errorCh)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error on creating cars",
			"error":   err,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cars succesfully created",
		"cars":    cars,
	})
}
