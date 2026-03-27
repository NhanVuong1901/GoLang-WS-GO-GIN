package main

import (
	"fmt"

	"ws/src/auth"
	"ws/src/chat"
	"ws/src/common"
	"ws/src/friend"
	"ws/src/notify"
	"ws/src/room"
	"ws/src/user"

	"github.com/gin-gonic/gin"
)

func main() {
	common.LoadEnv()
	db := common.MongoConnect()
	userRepo := user.NewRepository(db)
	friendRepo := friend.NewRepository(db)
	roomRepo := room.NewRepository(db)

	userController := user.NewController(userRepo)
	authController := auth.NewController(userRepo)
	friendController := friend.NewController(friendRepo)
	roomController := room.NewController(roomRepo)

	go chat.WS.Run()
	go notify.NotifyWS.Run()

	r := gin.Default()

	r.GET("/", func(ctx *gin.Context) {
		ctx.String(200, "Welcome to my chat server version 0.0.0.1")
	})
	r.POST("/api/register", userController.Register)
	r.POST("/api/login", authController.Login)
	r.POST("/api/friend/request", auth.JWTMiddleware(), friendController.SendRequest)
	r.POST("/api/friend/accept", auth.JWTMiddleware(), friendController.AcceptRequest)
	r.GET("/api/friends", auth.JWTMiddleware(), friendController.ListFriends)
	r.POST("/api/friend/refuse", auth.JWTMiddleware(), friendController.RefuseRequest)
	r.POST("/api/room", auth.JWTMiddleware(), roomController.Create)
	r.PATCH("/api/user", auth.JWTMiddleware(), userController.Update)

	r.GET("/ws", chat.ServerWS)
	r.GET("/ws/notify", notify.ServerWS)

	port := common.GetEnv("PORT")
	fmt.Println("Server is running at http://localhost" + port)
	r.Run(port)
}
