package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/emohankrishna/RMS/database"
	"github.com/emohankrishna/RMS/models"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")

func GetMenus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		result, err := menuCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while occurred while listing teh menu item"})
			return
		}
		var allMenus []bson.M
		if err = result.All(ctx, &allMenus); err != nil {
			color.Red("Error while decoding", err.Error())
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allMenus)
	}
}

func GetMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var menu models.Menu
		menuId := c.Param("menu_id")
		err := menuCollection.FindOne(ctx, bson.M{"menu_id": menuId}).Decode(&menu)
		if err != nil {
			color.Red("Failed to fetch menu in GetMenu")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while fetching the menu"})
			return
		}
		c.JSON(http.StatusOK, menu)
	}
}

func CreateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		var menu models.Menu
		if err := c.BindJSON(&menu); err != nil {
			color.Red("Not a valid request")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		menu.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.Menu_id = primitive.NewObjectID().Hex()
		validationErr := validate.Struct(menu)
		if validationErr != nil {
			color.Red("Not a valid menu")
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		insertResult, err := menuCollection.InsertOne(ctx, menu)
		if err != nil {
			color.Red("Error occurred while inserting Menu")
			c.JSON(http.StatusInternalServerError, gin.H{"error": validationErr.Error()})
			return
		}
		c.JSON(http.StatusOK, insertResult)
	}
}
func inTimeSpan(start, end, instant time.Time) bool {
	// current time should be less than start and end time
	return start.After(instant) && end.After(start)
}
func UpdateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		var menu models.Menu
		if err := c.BindJSON(&menu); err != nil {
			color.Red("Not a valid request")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var updateMenuObject primitive.D
		if menu.Start_Date != nil && menu.End_Date != nil {
			if inTimeSpan(*menu.Start_Date, *menu.End_Date, time.Now()) {
				msg := "kindly retype the time"
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
				return
			}
			var menu_id string = c.Param("menu_id")
			updateMenuObject = append(updateMenuObject, primitive.E{"start_date", menu.Start_Date})
			updateMenuObject = append(updateMenuObject, primitive.E{"end_date", menu.End_Date})
			if menu.Name != "" {
				updateMenuObject = append(updateMenuObject, bson.E{"name", menu.Name})
			}
			if menu.Category != "" {
				updateMenuObject = append(updateMenuObject, bson.E{"category", menu.Category})
			}

			menu.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			updateMenuObject = append(updateMenuObject, primitive.E{"updated_at", menu.Updated_at})
			filter := primitive.M{"menu_id": menu_id}
			upsert := true
			opt := options.UpdateOptions{
				Upsert: &upsert,
			}
			updateResult, err := menuCollection.UpdateOne(ctx, filter, bson.D{
				{"$set", updateMenuObject},
			}, &opt)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": updateResult})
			return

		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Start Date or End Date is invalid"})
	}
}
