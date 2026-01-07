package main

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

const OrderURL = "http://localhost:8082/orders"
const MenuURL = "http://localhost:8081/drinks"

func main() {
	router := gin.Default()

	// Маршруты к Menu Service
	router.GET("/api/drinks", getDrinks)
	router.GET("/api/drinks/:id", getDrinkByID)
	router.POST("/api/drinks", createDrink)

	// Маршруты к Orders Service
	router.GET("/api/orders", getOrders)
	router.GET("/api/orders/:id", getOrderByID)
	router.POST("/api/orders", createOrder)

	router.Run(":8080")
}

// ========== Menu Service ==========
func getDrinks(ctx *gin.Context) {
	resp, err := http.Get(MenuURL)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": "Menu Service недоступен"})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	ctx.Data(resp.StatusCode, "application/json", body)
}

func getDrinkByID(ctx *gin.Context) {
	id := ctx.Param("id")

	resp, err := http.Get(MenuURL + id)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": "Menu Service недоступен"})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	ctx.Data(resp.StatusCode, "application/json", body)
}

func createDrink(ctx *gin.Context) {
	// Читаем тело запроса от клиента
	reqBody, _ := io.ReadAll(ctx.Request.Body)

	// Пересылаем в Menu Service
	resp, err := http.Post(
		MenuURL,
		"application/json",
		bytes.NewReader(reqBody),
	)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": "Menu Service недоступен"})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	ctx.Data(resp.StatusCode, "application/json", body)
}

// ========== Orders Service ==========
func getOrders(ctx *gin.Context) {
	resp, err := http.Get(OrderURL)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": "Orders Service недоступен"})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	ctx.Data(resp.StatusCode, "application/json", body)
}

func getOrderByID(ctx *gin.Context) {
	id := ctx.Param("id")
	resp, err := http.Get(OrderURL + id)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": "Orders Service недоступен"})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	ctx.Data(resp.StatusCode, "application/json", body)
}

func createOrder(ctx *gin.Context) {
	// Читаем тело запроса от клиента
	reqBody, _ := io.ReadAll(ctx.Request.Body)

	// Пересылаем в Orders Service
	resp, err := http.Post(OrderURL, "application/json", bytes.NewReader(reqBody))
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": "Orders Service недоступен"})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	ctx.Data(resp.StatusCode, "application/json", body)
}
