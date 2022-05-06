package user

import (
	db "example/web-service-gin/db/sqlc"
	"example/web-service-gin/helper"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	db := c.MustGet("querier").(db.Querier)
	users, err := db.ListUser(c.Request.Context())
	if err != nil {
		log.Println("db.ListUser failed with:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": users})
}

func CreateUser(c *gin.Context) {
	var user db.CreateUserParams
	err := c.BindJSON(&user)
	if err != nil {
		log.Println("CreateUser: BindJSON Failed \n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	if !helper.CheckName(user.Firstname) || !helper.CheckName(user.Lastname) {
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

	db := c.MustGet("querier").(db.Querier)
	err = db.CreateUser(c.Request.Context(), user)

	if err != nil {
		message := fmt.Sprintf(
			"db.CreateUser failed with params %v \n",
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
	db := c.MustGet("querier").(db.Querier)
	err := db.DeactivateUser(c.Request.Context(), phone)

	if err != nil {
		message := fmt.Sprintf(
			"db.DeactivateUser failed with params %v \n",
			phone,
		)

		log.Println(message, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, "User Deactivated")
}
