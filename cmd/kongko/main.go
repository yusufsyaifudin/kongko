package main

import (
	"flag"
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin"
	"github.com/yusufsyaifudin/kongko/handler"
	"github.com/yusufsyaifudin/kongko/handler/kongko"
	"github.com/yusufsyaifudin/kongko/repo"
	"log"
	"net/http"
)

var addr = flag.String("addr", "localhost:9000", "http service address")
var wsAddr = flag.String("ws-addr", "ws://localhost:8083/mqtt", "websocket server address")

func main() {
	flag.Parse()
	//log.SetFlags(0)

	clientOptions := mqtt.NewClientOptions().AddBroker(*wsAddr).SetCleanSession(true)

	publisherClient := mqtt.NewClient(clientOptions)
	if client := publisherClient.Connect(); client.Wait() {
		err := client.Error()
		if err != nil {
			log.Printf("Failed creating mqtt connection. Error: %v", err)
		}
	}

	//mqttClient.Publish(topic, qos, retainLastMessage, payload)

	dataRepo := repo.NewMemDB() // use memdb for store data
	userHandler := kongko.NewUserHandler(dataRepo)
	chatHandler := kongko.NewChatHandler(dataRepo, publisherClient)

	server := gin.Default()
	gin.SetMode("debug")

	publicResource := server.Group("/api/v1")
	publicResource.POST("/register", userHandler.Register)
	publicResource.POST("/login", userHandler.Login)

	protectedResource := server.Group("/api/v1")
	protectedResource.Use(handler.ProtectedResource(dataRepo))
	protectedResource.GET("/chats", chatHandler.ListConversation)
	protectedResource.POST("/chats", chatHandler.CreateConversation)
	protectedResource.GET("/messages", chatHandler.GetMessageInRoomId)
	protectedResource.POST("/messages", chatHandler.PostMessage)

	server.StaticFS("/home", http.Dir("web"))

	err := server.Run(*addr)
	if err != nil {
		log.Fatal(err)
	}
}
