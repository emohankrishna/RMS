package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/emohankrishna/RMS/database"
	"github.com/emohankrishna/RMS/helpers"
	"github.com/emohankrishna/RMS/models"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		var user models.User
		if err := c.BindJSON(&user); err != nil {
			color.Red("Invalid Request while signing up")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		validationErr := validate.Struct(user)
		if validationErr != nil {
			color.Red("Validation Error while signing up")
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		user.Created_at = time.Now()
		user.Updated_at = time.Now()
		user.Password = helpers.GetHash([]byte(user.Password))

		userClaims := helpers.UserClaims{
			User_id: user.Email,
		}
		tokenString, tokenErr := helpers.CreateToken(&userClaims)
		if tokenErr != nil {
			color.Red("Error While creating the token")
			c.JSON(http.StatusInternalServerError, gin.H{"error": tokenErr.Error()})
		}
		user.Token = tokenString
		user.RefreshToken = tokenString
		inserted, err := userCollection.InsertOne(ctx, user)
		if err != nil {
			color.Red("Error while Inserting the user document")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, inserted)

	}
}
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		var loginUser models.LoginUser
		if err := c.BindJSON(&loginUser); err != nil {
			color.Red("Invalid Request while login")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		validationErr := validate.Struct(loginUser)
		if validationErr != nil {
			color.Red("Validation Error while signing up")
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		var foundUser models.User
		err := userCollection.FindOne(ctx, bson.M{"email": loginUser.Email}).Decode(&foundUser)
		if err != nil {
			color.Red("Error while Finding the user document")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var isAuthenticated bool = helpers.VerifyPassword([]byte(foundUser.Password), []byte(loginUser.Password))
		if !isAuthenticated {
			color.Red("Not a valid User")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, foundUser)
	}
}
