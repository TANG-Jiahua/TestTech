package main

import (
	"TestTech/Service"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"

	//"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/gin-gonic/gin"
	"net/http"
)

//type Collection struct {
//	DB *mongo.Collection
//}

func initDB(){
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	//defer func() {
	//	if err = client.Disconnect(ctx); err != nil {
	//		panic(err)
	//	}
	//}()
	Service.MyUser = &Service.UserModel{
		DB: client.Database("test").Collection("user"),
	}
}

func main ()  {

	initDB()

	r:=gin.Default()
	r.GET("/", func(c  *gin.Context) {c.String(http.StatusOK,"ok")})

	r.POST("/add/users",Service.AddUser)
	r.POST("/login",Service.Login)
	r.DELETE("/delete/user/:id",Service.DeleteUser)
	r.GET("users/list",Service.GetUserList)
	r.GET("/user/:id",Service.GetUserById)
	r.PUT("/user/:id",Service.UpdateUser)
	r.Run()
}
