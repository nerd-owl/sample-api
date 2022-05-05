package user

import (
	"example/web-service-gin/handler/pgx"
	"example/web-service-gin/helper"
	"fmt"

	"log"
	"net/http"

	"github.com/gin-gonic/gin"
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

func GetUsers(c *gin.Context) {
	conn := c.MustGet("databaseConn").(pgx.DbConn)
	selectQuery := "select * from kuser"
	rows, err := conn.Query(c.Request.Context(), selectQuery)
	if err != nil {
		message := fmt.Sprintf("GetUsers: %v failed \n", selectQuery)
		log.Println(message, err)
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
			log.Println(
				"GetUsers: rows.Scan Failed \n",
				err,
			)
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
		log.Println("CreateUser: BindJSON Failed \n", err)
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

	conn := c.MustGet("databaseConn").(pgx.DbConn)
	insertQuery := `INSERT INTO kuser (FirstName, LastName, Phone, Addr)
	VALUES ($1, $2, $3, $4)`
	_, err = conn.Exec(
		c.Request.Context(),
		insertQuery,
		user.Fname,
		user.Lname,
		user.Phone,
		user.Address,
	)

	if err != nil {
		message := fmt.Sprintf(
			"CreateUser: %v failed with params %v \n",
			insertQuery,
			user,
		)

		log.Println(message, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

func DeactivateUser(c *gin.Context) {
	phone := c.Param("phone")
	conn := c.MustGet("databaseConn").(pgx.DbConn)
	updateQuery := `UPDATE kuser
	SET Active = False
	WHERE Phone = $1`
	_, err := conn.Exec(c.Request.Context(), updateQuery, phone)

	if err != nil {
		message := fmt.Sprintf(
			"DeactivateUser: %v failed with params %v \n",
			updateQuery,
			phone,
		)

		log.Println(message, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, "User Deactivated")
}
