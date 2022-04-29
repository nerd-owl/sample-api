package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
)

type User struct {
	UserWithoutActivity
	Active bool
}

type UserWithoutActivity struct {
	Fname   string
	Lname   string
	Phone   string
	Address string
}

func checkPhone(phone string) bool {
	validPhone := regexp.MustCompile("^[0-9]{10}$")
	return validPhone.MatchString(phone)
}

func checkName(name string) bool {
	validName := regexp.MustCompile("^[a-zA-Z]{3,}$")
	return validName.MatchString(name)
}

func attachPgConn(conn *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("databaseConn", conn)
		c.Next()
	}
}

func getUsers(c *gin.Context) {
	conn := c.MustGet("databaseConn").(*pgx.Conn)

	rows, err := conn.Query(context.Background(), "select * from kuser")
	if err != nil {
		log.Println("Query failed", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	var users []User

	for rows.Next() {
		var user User = User{}
		err := rows.Scan(
			&user.Fname,
			&user.Lname,
			&user.Phone,
			&user.Address,
			&user.Active,
		)
		if err != nil {
			log.Println("row.Scan Failed", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		users = append(users, user)
	}

	c.JSON(http.StatusOK, gin.H{"data": users})
}

func createUser(c *gin.Context) {
	var user UserWithoutActivity
	err := c.BindJSON(&user)
	if err != nil {
		log.Println("BindJSON Failed", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	if !checkName(user.Fname) || !checkName(user.Lname) {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "FirstName and LastName should not have spaces and must only contain letters",
			},
		)
		return
	}

	if !checkPhone(user.Phone) {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "Phone must be exactly 10 digits and only numbers",
			},
		)
		return
	}

	conn := c.MustGet("databaseConn").(*pgx.Conn)
	insertQuery := `INSERT INTO kuser (FirstName, LastName, Phone, Addr)
	VALUES ($1, $2, $3, $4)`
	_, err = conn.Exec(
		context.Background(),
		insertQuery,
		user.Fname,
		user.Lname,
		user.Phone,
		user.Address,
	)

	if err != nil {
		log.Println("BindJSON Failed", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

func deactivateUser(c *gin.Context) {
	phone := c.Param("phone")
	conn := c.MustGet("databaseConn").(*pgx.Conn)
	updateQuery := `UPDATE kuser
	SET Active = False
	WHERE Phone = $1`
	_, err := conn.Exec(context.Background(), updateQuery, phone)

	if err != nil {
		log.Println("Update Query Failed", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, "User Deactivated")
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

	router.GET("/users", getUsers)
	router.POST("/users", createUser)
	router.GET("/deactivate/:phone", deactivateUser)
	router.Run(fmt.Sprintf("localhost:%v", os.Getenv("HTTP_PORT")))
}
