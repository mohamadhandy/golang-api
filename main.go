package main

import (
	"fmt"
	"golang-api/auth"
	"golang-api/handler"
	"golang-api/helper"
	"golang-api/user"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
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
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar)
	router.Run()

	// buat layering

	// input
	// handler mapping input dari user ke struct
	// service mapping ke struct user
	// repository save struct user ke db
	// db
}

func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		// Bearer token
		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}
		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		claim, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		userId := int(claim["user_id"].(float64))
		user, err := userService.GetUserByID(userId)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		c.Set("currentUser", user)
	}
}

// midleware:
// 1. ambil nilai header authorization: Bearer token
// 2. dari header authorization, kita ambil nilai tokennya saja
// 3. kita validasi token(pakai service ValidateToken)
// 4. Jika valid tokennya, kita ambil user_id
// 5. Ambil user dari db berdasarkan user_id lewat service.
// 6. kita set context isinya user
// context: sebuah tempat untuk menyimpan suatu nilai yang pada akhirnya bisa di get/diambil dari tempat yang lain.
