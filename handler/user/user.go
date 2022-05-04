package user

import (
	"example/web-service-gin/helper"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
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
	conn := c.MustGet("databaseConn").(*pgx.Conn)
	selectQuery := "select * from kuser"
	rows, err := conn.Query(c.Request.Context(), selectQuery)
	if err != nil {
		log.Println("GetUsers: Query failed for: ", selectQuery, err)
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
				"GetUsers: row.Scan Failed for: ",
				selectQuery,
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

	conn := c.MustGet("databaseConn").(*pgx.Conn)
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
		log.Println("CreateUser: ", insertQuery, "failed \n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

func DeactivateUser(c *gin.Context) {
	phone := c.Param("phone")
	conn := c.MustGet("databaseConn").(*pgx.Conn)
	updateQuery := `UPDATE kuser
	SET Active = False
	WHERE Phone = $1`
	_, err := conn.Exec(c.Request.Context(), updateQuery, phone)

	if err != nil {
		log.Println("DeactivateUser: ", updateQuery, "failed. \n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, "User Deactivated")
}
