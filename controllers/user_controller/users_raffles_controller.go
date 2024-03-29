package usercont

import (
	"context"
	rafcont "go_rest_api_skeleton/controllers/raffle_contoller"
	"go_rest_api_skeleton/models"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//! Kullanıcının kayıtlı olduğu çekilişleri getiren istek.
func SelectUsersRaffles(c echo.Context) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var result models.UserModel
	//! doğrulanan tokenın içerisindeki uid parametresini alır string e dönüştürür
	uid := c.Param("uid")

	gt, errInt := strconv.Atoi(c.QueryParam("gt"))
	if errInt != nil {
		panic(errInt)
	}

	err := userCollection.FindOne(ctx, bson.D{{Key: "uid", Value: uid}}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.JSON(http.StatusBadRequest, models.Response{Body: &echo.Map{"error": err.Error()}})
		}
		panic(err)
	}

	var raffleIdList []primitive.ObjectID
	//! kullanıcının kayıt olduğu raffleları listeye atar.
	for _, value := range result.SubscribedRaffles {
		raffleIdList = append(raffleIdList, value.RaffleId)
	}
	if len(raffleIdList) == 0 {
		return c.JSON(http.StatusNoContent, models.Response{Body: &echo.Map{"data": "There is no any raffle"}})
	}
	var usersRaffelList models.UsersRaffleList
	// bunun boşlupunu kontrol et
	usersRaffelList.RaffleList = rafcont.GetUsersRaffles(gt, raffleIdList)

	return c.JSON(http.StatusOK, models.Response{Body: &echo.Map{"data": usersRaffelList}})
}
