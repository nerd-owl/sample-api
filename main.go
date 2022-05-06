package main

import (
	"database/sql"
	"log"

	db "example/web-service-gin/db/sqlc"
	"example/web-service-gin/handler/user"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func attachPgConn(querier db.Querier) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("querier", querier)
		c.Next()
	}
}

const (
	dbDriver = "postgres"
	dbSource = "postgres://postgres:mysecretpassword@localhost:5432/postgres?sslmode=disable"
)

func main() {

	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	defer conn.Close()

	querier := db.New(conn)

	router := gin.Default()
	router.Use(attachPgConn(querier))

	router.GET("/users", user.GetUsers)
	router.POST("/users", user.CreateUser)
	router.PUT("/deactivate/:phone", user.DeactivateUser)
	router.Run("localhost:4030")
}
