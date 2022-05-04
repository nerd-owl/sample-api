package user

import (
	"context"
	"example/web-service-gin/helper"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct {
	Fname   string
	Lname   string
	Phone   string
	Address string
	Active  bool
}

type UserWithoutActivity struct {
	Fname   string
	Lname   string
	Phone   string
	Address string
}

type Rows interface {
	Next() bool
	Scan(dest ...interface{}) error
}

type dbconn interface {
	Query(
		ctx context.Context,
		sql string,
		args ...interface{},
	) (Rows, error)
	Exec(
		ctx context.Context,
		sql string,
		arguments ...interface{},
	) (interface{}, error)
}

func GetUsers(c *gin.Context) {
	conn := c.MustGet("databaseConn").(dbconn)

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

func CreateUser(c *gin.Context) {
	var user UserWithoutActivity
	err := c.BindJSON(&user)
	if err != nil {
		log.Println("BindJSON Failed", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	if !helper.CheckName(user.Fname) || !helper.CheckName(user.Lname) {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "FirstName and LastName should not have spaces and must only contain letters",
			},
		)
		return
	}

	if !helper.CheckPhone(user.Phone) {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "Phone must be exactly 10 digits and only numbers",
			},
		)
		return
	}

	conn := c.MustGet("databaseConn").(dbconn)
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

func DeactivateUser(c *gin.Context) {
	phone := c.Param("phone")
	conn := c.MustGet("databaseConn").(dbconn)
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
