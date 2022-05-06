package user

import (
	db "example/web-service-gin/db/sqlc"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	store := c.MustGet("store").(db.Querier)
	users, err := store.ListUser(c.Request.Context())
	if err != nil {
		log.Println("store.ListUser failed with:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": users})
}

type CreateUserRequest struct {
	Firstname string `json:"firstname" binding:"required,alpha"`
	Lastname  string `json:"lastname" binding:"required,alpha"`
	Phone     string `json:"phone" binding:"required,number,len=10"`
	Addr      string `json:"addr" binding:"required"`
}

func CreateUser(c *gin.Context) {
	var user CreateUserRequest
	err := c.ShouldBindJSON(&user)
	if err != nil {
		log.Println("CreateUser: BindJSON Failed \n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	newUser := db.CreateUserParams{
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Phone:     user.Phone,
		Addr:      user.Addr,
	}

	store := c.MustGet("store").(db.Querier)
	err = store.CreateUser(c.Request.Context(), newUser)

	if err != nil {
		message := fmt.Sprintf(
			"store.CreateUser failed with params %v \n",
			newUser,
		)

		log.Println(message, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

func DeactivateUser(c *gin.Context) {
	phone := c.Param("phone")
	store := c.MustGet("store").(db.Querier)
	err := store.DeactivateUser(c.Request.Context(), phone)

	if err != nil {
		message := fmt.Sprintf(
			"store.DeactivateUser failed with params %v \n",
			phone,
		)

		log.Println(message, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, "User Deactivated")
}
