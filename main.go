package main

import (
	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

type Message struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

var db *gorm.DB

func initDB() {
	dsn := "host=localhost user=postgres password=yourpassword dbname=postgres port=5433 sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	db.AutoMigrate(&Message{})
}

func GetHandler(c echo.Context) error {
	var message []Message

	if err := db.Find(&message).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Status:  "error",
			Message: "Could not find message",
		})
	}

	return c.JSON(http.StatusOK, &message)
}

func PostHandler(c echo.Context) error {
	var message Message

	if err := c.Bind(&message); err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Could not add message",
		})
	}

	if err := db.Create(&message).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Status:  "error",
			Message: "Could not create message",
		})
	}

	return c.JSON(http.StatusOK, Response{
		Status:  "success",
		Message: "Message added successfully",
	})
}

func UpdateHandler(c echo.Context) error {
	var updateMessage Message
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Invalid id",
		})
	}

	if err := c.Bind(&updateMessage); err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Could not add message",
		})
	}

	if err := db.Model(&Message{}).Where("id = ?", id).Update("text", updateMessage.Text).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Status:  "error",
			Message: "Could not update message",
		})
	}
	return c.JSON(http.StatusOK, Response{
		Status:  "success",
		Message: "Message updated successfully",
	})
}

func DeleteHandler(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Invalid id",
		})
	}

	if err := db.Delete(&Message{}, id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Status:  "error",
			Message: "Could not delete message",
		})
	}

	return c.JSON(http.StatusOK, Response{
		Status:  "success",
		Message: "Message deleted successfully",
	})
}

func main() {
	initDB()
	e := echo.New()
	e.GET("/messages", GetHandler)
	e.POST("/messages", PostHandler)
	e.PATCH("/messages/:id", UpdateHandler)
	e.DELETE("/messages/:id", DeleteHandler)
	e.Start(":8080")
}
