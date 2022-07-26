package main

import (
	"fmt"
	"golang-api/auth"
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
	fmt.Println("db", &db)
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	authService := auth.NewService()

	token, err := authService.ValidateToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo1fQ.8cfLZ-lWx-zb7bP-dgSh03OoI4dafd3eHpXxVYEglNo")
	if err != nil {
		fmt.Println("ERROR!!!")
	}
	if token.Valid {
		fmt.Println("VALID TOKEN!!!")
	} else {
		fmt.Println("INVALID TOKEN!!!")
	}

	// authService.GenerateToken(1001)

	// fmt.Println("userRepository", &userRepository)
	// fmt.Println("userService", &userService)

	userHandler := handler.NewUserHandler(userService, authService)
	router := gin.Default()

	api := router.Group("/api/v1")
	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)
	api.POST("/avatars", userHandler.UploadAvatar)
	router.Run()

	// buat layering

	// input
	// handler mapping input dari user ke struct
	// service mapping ke struct user
	// repository save struct user ke db
	// db
}
