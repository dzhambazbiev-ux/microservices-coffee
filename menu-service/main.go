package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Drink struct {
	gorm.Model
	Name    string  `json:"name"`
	Price   float64 `json:"price"`
	InStock bool    `json:"in_stock"`
}

type CreateDrinkRequest struct {
	Name    string  `json:"name"`
	Price   float64 `json:"price"`
	InStock bool    `json:"in_stock"`
}

var db *gorm.DB

func main() {
	dsn := "host=localhost user=postgres password=postgres dbname=menu_db port=5432 sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("не удалось подключиться к бд" + err.Error())
	}

	db.AutoMigrate(&Drink{})

	router := gin.Default()

	router.GET("/drinks", listDrinks)
	router.GET("/drinks/:id", getDrink)
	router.POST("/drinks", createDrink)

	router.Run(":8081")
}

func listDrinks(ctx *gin.Context) {
	var drinks []Drink
	db.Find(&drinks)
	ctx.JSON(http.StatusOK, drinks)
}

func getDrink(ctx *gin.Context) {
	id := ctx.Param("id")

	var drink Drink
	if err := db.First(&drink, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "напиток не найден"})
		return
	}

	ctx.JSON(http.StatusOK, drink)
}

func createDrink(ctx *gin.Context) {
	var req CreateDrinkRequest
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "некорректный JSON"})
		return
	}

	drink := Drink{
		Name:    req.Name,
		Price:   req.Price,
		InStock: req.InStock,
	}

	if err := db.Create(&drink).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, drink)
}
