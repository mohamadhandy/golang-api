package main

import (
	"fmt"
	"golang-api/handler"
	"golang-api/user"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:@tcp(127.0.0.1:3306)/bwastartup?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	// userInput := user.RegisterUserInput{}
	// userInput.Name = "service 5"
	// userInput.Email = "service5@gmail.com"
	// userInput.Occupation = "Admin test"
	// userInput.Password = "password"

	// userService.RegisterUser(userInput)
	fmt.Println("userRepository", &userRepository)
	fmt.Println("userService", &userService)

	userHandler := handler.NewUserHandler(userService)
	router := gin.Default()

	api := router.Group("/api/v1")
	api.POST("/users", userHandler.RegisterUser)
	router.Run()

	// buat layering

	// input
	// handler mapping input dari user ke struct
	// service mapping ke struct user
	// repository save struct user ke db
	// db
}
