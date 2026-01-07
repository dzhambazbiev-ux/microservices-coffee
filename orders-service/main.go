package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	DrinkID   uint    `json:"drink_id"`
	DrinkName string  `json:"drink_name"` // Копируем из Menu Service
	Price     float64 `json:"price"`      // Копируем из Menu Service
	Quantity  int     `json:"quantity"`
	Total     float64 `json:"total"`
	Status    string  `json:"status"` // pending, completed, cancelled
}

type CreateOrderRequest struct {
	DrinkID  uint `json:"drink_id"`
	Quantity int  `json:"quantity"`
}

// Структура для ответа от Menu Service
type DrinkFromMenu struct {
	ID      uint    `json:"id"`
	Name    string  `json:"name"`
	Price   float64 `json:"price"`
	InStock bool    `json:"in_stock"`
}

var db *gorm.DB

// Адрес Menu Service
const menuServiceURL = "http://localhost:8081"

func main() {
	dsn := "host=localhost user=postgres password=postgres dbname=orders_db port=5432 sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Ошибка при подключении к бд" + err.Error())
	}

	db.AutoMigrate(&Order{})

	router := gin.Default()

	router.GET("/orders", listOrders)
	router.GET("/orders/:id", getOrder)
	router.POST("/orders", createOrder)

	router.Run(":8082")
}

func createOrder(ctx *gin.Context) {
	var req CreateOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "некорректный JSON"})
		return
	}

	if req.Quantity <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Количество должно быть больше 0"})
		return
	}

	// ====== ГЛАВНАЯ СВЯЗКА: запрос к Menu Service ======
	drink, err := getDrinkFromMenuService(req.DrinkID)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": "Menu Service недоступен" + err.Error()})
		return
	}

	if drink == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "напиток не найден"})
		return
	}

	if !drink.InStock {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "напитка нет в наличии"})
		return
	}
	// ====================================================

	order := Order{
		DrinkID:   drink.ID,
		DrinkName: drink.Name,
		Price:     drink.Price,
		Quantity:  req.Quantity,
		Total:     drink.Price * float64(req.Quantity),
		Status:    "pending",
	}

	db.Create(&order)
	ctx.JSON(http.StatusOK, order)
}

func getDrinkFromMenuService(drinkID uint) (*DrinkFromMenu, error) {
	url := fmt.Sprintf("%s/drinks/%d", menuServiceURL, drinkID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	// Напиток не найден
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	// Другая ошибка
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Menu Service вернул статус: %d", resp.StatusCode)
	}

	// Читаем и парсим ответ
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var drink DrinkFromMenu
	if err := json.Unmarshal(body, &drink); err != nil {
		return nil, err
	}

	return &drink, nil
}

func getOrder(ctx *gin.Context) {
	id := ctx.Param("id")

	var order Order
	if err := db.First(&order, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "заказ не найден"})
		return
	}

	ctx.JSON(http.StatusOK, order)
}

func listOrders(ctx *gin.Context) {
	var orders []Order
	if err := db.Find(&orders).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, orders)
}
