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
<<<<<<< HEAD
		AllowOrigins: []string{"https://waysbeans-99.vercel.app/", "https://waysbeans-fe-git-main-kyuba00.vercel.app/", "https://waysbeans-f0859gxhq-kyuba00.vercel.app/"},
=======
		AllowOrigins: []string{"*", "https://waysbeans-99.vercel.app"},
>>>>>>> 6fec971ba5aa83d5e2cd5a9d609517f80aa94b6e
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
