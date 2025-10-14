package main

import (
    "context"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type Borrow struct {
    UserID string `json:"user_id"`
    BookID string `json:"book_id"`
}

var collection *mongo.Collection

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    client, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://mongo-service:27017"))
    collection = client.Database("borrows").Collection("entries")

    r := gin.Default()

    r.GET("/borrows", func(c *gin.Context) {
        var results []Borrow
        cursor, _ := collection.Find(ctx, bson.M{})
        cursor.All(ctx, &results)
        c.JSON(http.StatusOK, results)
    })

    r.POST("/borrows", func(c *gin.Context) {
        var borrow Borrow
        if err := c.ShouldBindJSON(&borrow); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        collection.InsertOne(ctx, borrow)
        c.JSON(http.StatusCreated, borrow)
    })

    r.Run(":3002")
}
