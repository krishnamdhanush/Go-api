package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/krishnamdhanush/web-dawn/configs"
	"go.mongodb.org/mongo-driver/bson"
)

type Album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// var albums = []Album{
// 	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
// 	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
// 	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
// }

// getAlbums responds with the list of all albums as JSON.
func GetAlbums(c *gin.Context) {
	ctx := context.Background()
	albumsResult := []Album{}
	albumCollection := configs.GetCollection(configs.DB, "albums")

	albums_cursor, err := albumCollection.Find(ctx, bson.D{})
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, nil)
		return
	}
	for albums_cursor.Next(context.TODO()) {
		var result Album
		if err := albums_cursor.Decode(&result); err != nil {
			log.Fatal(err)
		}
		albumsResult = append(albumsResult, result)
	}
	if err := albums_cursor.Err(); err != nil {
		log.Fatal(err)
		c.IndentedJSON(http.StatusInternalServerError, nil)
		return
	}

	//push to redis
	err = redisInsert("albums:all", albumsResult)
	if err != nil {
		log.Fatalf("couldn't save to redis: %v", err)
	}

	c.IndentedJSON(http.StatusOK, albumsResult)
}

func PostAlbum(c *gin.Context) {
	ctx := context.Background()
	var newAlbum Album
	if err := c.BindJSON(&newAlbum); err != nil {
		log.Fatal(err)
		c.IndentedJSON(http.StatusBadRequest, "Invalid album object")
		return
	}
	albumCollection := configs.GetCollection(configs.DB, "albums")

	_, err := albumCollection.InsertOne(ctx, newAlbum)
	if err != nil {
		log.Fatal(err)
		c.IndentedJSON(http.StatusInternalServerError, "Couldn't save album. Try again!")
		return
	}

	//push to redis
	err = redisInsert(fmt.Sprintf("albums:%s", newAlbum.ID), newAlbum)
	if err != nil {
		log.Fatalf("couldn't save to redis: %v", err)
	}
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

func GetAlbumById(c *gin.Context) {
	id, check := c.Params.Get("id")
	if !check {
		panic("No id param found")
	}
	var result Album
	albumCollection := configs.GetCollection(configs.DB, "albums")
	ctx := context.Background()

	err := albumCollection.FindOne(ctx, bson.D{{Key: "id", Value: id}}).Decode(&result)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, "No album with such Id")
		return
	}

	//push to redis
	err = redisInsert(fmt.Sprintf("albums:%s", id), result)
	if err != nil {
		log.Fatalf("couldn't save to redis: %v", err)
	}

	c.IndentedJSON(http.StatusCreated, result)
}

func redisInsert(key string, val interface{}) error {
	ctx := context.Background()
	marshalledObject, err := json.Marshal(val)
	if err != nil {
		return err
	}
	err = configs.RedisClient.Set(ctx, key, marshalledObject, time.Second*10).Err()
	if err != nil {
		return err
	}
	return nil
}
