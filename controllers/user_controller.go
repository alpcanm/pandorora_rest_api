package controllers

import (
	"context"
	config "go_rest_api_skeleton/config"
	"go_rest_api_skeleton/models"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var collection *mongo.Collection = config.GetCollection(config.DB, "users")
var validate = validator.New()

func InsertAUser(c echo.Context) error {
	//10 saniye süreli bir context oluşturur.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//tüm işlemler bitince contexti iptal eder.
	defer cancel()
	var user models.UserModel
	// request'in içindeki verileri alır ve user model'e yerleştirir.
	if err := c.Bind(&user); err != nil {
		//bir hata varsa döner
		return c.JSON(http.StatusBadRequest, models.Response{Body: &echo.Map{"error": err.Error()}})
	}
	// requestin içindeki veriyi validate eder.
	if err := validate.Struct(user); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{Body: &echo.Map{"error": err.Error()}})
	}
	//user in uid sine mongodb _uid yerleştirid
	user.Uid = primitive.NewObjectID()
	//veritabanına ekler.
	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{Body: &echo.Map{"error": err.Error()}})
	}
	//sonuç döner
	return c.JSON(http.StatusCreated, models.Response{Body: &echo.Map{"data": result}})
}

func SelectAUser(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var result bson.M
	// parametreyi uid değişkenine atar
	uid := c.Param("uid")
	// filtreye göre istenilen veriyi getirir ve result değişkenine atar.
	err := collection.FindOne(ctx, models.UserFilterModel{Uid: uid}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.JSON(http.StatusBadRequest, models.Response{Body: &echo.Map{"error": err.Error()}})
		}
		panic(err)
	}
	return c.JSON(http.StatusCreated, models.Response{Body: &echo.Map{"data": result}})
}
