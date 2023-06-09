package controller

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"Backend/models"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateUser(c echo.Context) error {
	user := new(models.User)
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	
	client := models.MongoConnect()
	defer client.Disconnect(context.TODO())

	coll := models.MongoCollection("User", client)
	cursor, err := coll.Find(context.Background(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var userTemp models.User
		if err := cursor.Decode(&userTemp); err != nil {
			log.Fatal(err)
		}
		if(userTemp.UserName == user.UserName){
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Username sudah ada, silahkan cari username lain"})
		}
	}



	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var table = models.MongoCollection("User", client)

	var newUser = models.User{
		ID: primitive.NewObjectID(),
		UserName : user.UserName,
		Password:  user.Password,
	}
	_, errInsert := table.InsertOne(ctx, newUser)

	if errInsert != nil {
		return c.String(http.StatusInternalServerError, "Server Error")
	}
	return c.JSON(http.StatusOK, newUser)
}

func GetUser(c echo.Context) error {
	username := c.QueryParam("username")
	password := c.QueryParam("password")

	fmt.Println(username)
	fmt.Println(password)
	client := models.MongoConnect()
	defer client.Disconnect(context.TODO())

	coll := models.MongoCollection("User", client)

	cursor, err := coll.Find(context.Background(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var userTemp models.User
		if err := cursor.Decode(&userTemp); err != nil {
			log.Fatal(err)
		}
		if(userTemp.UserName == username){
			fmt.Println(userTemp)
			if(userTemp.Password == password){
				return c.JSON(http.StatusOK, userTemp)
			}else{
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid password"})
			}
		}
	}
	return c.JSON(http.StatusBadRequest, map[string]string{"error": "Username tidak ditemukan"})
}