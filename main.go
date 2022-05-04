package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"example/web-service-gin/handler/user"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
)

func attachPgConn(conn *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("databaseConn", conn)
		c.Next()
	}
}

func main() {

	err := godotenv.Load("./configs/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	router := gin.Default()
	// url := "postgres://postgres:mysecretpassword@localhost:5432/postgres"
	url := fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		log.Fatal("Unable to connect to database \n", err)
	}
	defer conn.Close(context.Background())

	router.Use(attachPgConn(conn))

	router.GET("/users", user.GetUsers)
	router.POST("/users", user.CreateUser)
	router.PUT("/deactivate/:phone", user.DeactivateUser)
	router.Run(fmt.Sprintf("localhost:%v", os.Getenv("HTTP_PORT")))
}
