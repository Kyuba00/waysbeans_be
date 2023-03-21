package main

import (
	"fmt"
	"os"
	"waysbeans_be/database"
	"waysbeans_be/pkg/mysql"
	"waysbeans_be/routes"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	errEnv := godotenv.Load()
	if errEnv != nil {
		panic("Failed to load env file")
	}

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*", "https://waysbeans-99.vercel.app"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PATCH, echo.DELETE},
		AllowHeaders: []string{"X-Requested-With", "Content-Type", "Authorization"},
	}))
	
	mysql.DatabaseInit()
	database.RunMigration()

	routes.RouteInit(e.Group("/api/v1"))

	e.Static("/uploads", "./uploads")

	PORT := os.Getenv("PORT")

	fmt.Println("server running localhost:" + PORT)
	e.Logger.Fatal(e.Start(":" + PORT))
}
