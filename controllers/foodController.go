package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/emohankrishna/RMS/database"
	"github.com/emohankrishna/RMS/models"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	validator "github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")

var validate = validator.New()

func GetFoods() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		recordsPerPage, err := strconv.Atoi(c.DefaultQuery("recordsPerPage", "10"))

		if err != nil || recordsPerPage < 1 {
			recordsPerPage = 10
		}
		page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
		if err != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordsPerPage
		// startIndex, err = strconv.Atoi(c.Query("startIndex"))

		matchStage := bson.D{
			{"$match", bson.D{{}}},
		}
		groupStage := bson.D{
			{"$group",
				bson.D{
					{"_id", primitive.Null{}},
					{"total_count", bson.D{{"$sum", 1}}},
					{"data", bson.D{{"$push", "$$ROOT"}}},
				},
			},
		}
		projectStage := bson.D{
			{"$project",
				bson.D{
					{"_id", 0},
					{"total_count", 1},
					{"data",
						bson.D{
							{"$slice",
								bson.A{
									"$data",
									startIndex,
									recordsPerPage,
								},
							},
						},
					},
				},
			},
		}
		color.Red("recordsPerPage%d startIndex %d", recordsPerPage, startIndex)
		var allFood []bson.M
		resultCursor, err := foodCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage,
		})
		if err != nil {
			color.Red("Failed to fetch All food ", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while fetching the foods"})
			return
		}
		if err := resultCursor.All(ctx, &allFood); err != nil {
			color.Red("Error while decoding", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while decoding the foods"})
			return
		}
		c.JSON(http.StatusOK, allFood)
	}
}

func GetFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		foodId := c.Param("food_id")
		var food models.Food

		err := foodCollection.FindOne(ctx, bson.M{"food_id": foodId}).Decode(&food)
		if err != nil {
			color.Red("Failed to fetch food in GetFood")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while fetching the food"})
			return
		}
		c.JSON(http.StatusOK, food)
	}
}

func CreateFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		var menu models.Menu
		var food models.Food
		if err := c.BindJSON(&food); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		food.Food_id = fmt.Sprintf("F-%s", primitive.NewObjectID().Hex())
		validationErr := validate.Struct(food)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		err := menuCollection.FindOne(ctx, bson.M{"menu_id": food.Menu_id}).Decode(&menu)
		if err != nil {
			color.Red("Menu not found")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Menu not found for the given menu_id"})
			return
		}
		food.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		food.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		result, insertErr := foodCollection.InsertOne(ctx, food)
		if insertErr != nil {
			msg := fmt.Sprintln("Error while creating food")
			color.Red(msg)
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}
func UpdateFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		var menu models.Menu
		var food models.Food
		if err := c.BindJSON(&food); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		validationErr := validate.Struct(food)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		err := menuCollection.FindOne(ctx, bson.M{"menu_id": food.Menu_id}).Decode(&menu)
		if err != nil {
			color.Red("Menu not found")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Menu not found for the given menu_id"})
			return
		}
		var updateFoodObject bson.D
		if food.Name != "" {
			updateFoodObject = append(updateFoodObject, bson.E{"name", food.Name})
		}
		if food.Price != nil {
			updateFoodObject = append(updateFoodObject, bson.E{"price", *food.Price})
		}
		if food.Food_image != "" {
			updateFoodObject = append(updateFoodObject, bson.E{"food_image", food.Food_image})
		}
		if food.Menu_id != "" {
			updateFoodObject = append(updateFoodObject, bson.E{"menu_id", food.Menu_id})
		}
		food.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateFoodObject = append(updateFoodObject, bson.E{"updated_at", food.Updated_at})
		foodId := c.Param("food_id")
		filter := bson.D{{"food_id", foodId}}
		update := bson.D{{"$set", updateFoodObject}}
		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		updateResult, err := foodCollection.UpdateOne(ctx, filter, update, &opt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"results": updateResult})
	}
}
