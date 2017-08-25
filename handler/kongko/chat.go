package kongko

import (
	"encoding/json"
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin"
	"github.com/yusufsyaifudin/kongko/model"
	"github.com/yusufsyaifudin/kongko/repo"
)

type ChatHandler struct {
	DataRepo  repo.DataRepository
	Publisher mqtt.Client
}

func NewChatHandler(dataRepo repo.DataRepository, publisher mqtt.Client) (chatHandler *ChatHandler) {
	chatHandler = &ChatHandler{
		DataRepo:  dataRepo,
		Publisher: publisher,
	}
	return
}

func (c *ChatHandler) ListConversation(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)

	chats, _ := c.DataRepo.GetConversationList(user)

	ctx.JSON(200, map[string]interface{}{
		"user": map[string]interface{}{
			"email": user.Email,
		},
		"chats": chats,
	})
}

func (c *ChatHandler) CreateConversation(ctx *gin.Context) {
	name := ctx.PostForm("name")
	emails := ctx.PostFormArray("emails")

	if name == "" {
		ctx.JSON(422, map[string]interface{}{
			"error": "name cannot be empty",
		})
		ctx.Abort()
		return
	}

	if len(emails) < 1 {
		ctx.JSON(422, map[string]interface{}{
			"error": "at least one email must be present",
		})
		ctx.Abort()
		return
	}

	users := []*model.User{}
	alreadyAdded := make(map[string]bool)
	user := ctx.MustGet("user").(*model.User)

	users = append(users, user)
	alreadyAdded[user.Email] = true

	for _, email := range emails {
		u, err := c.DataRepo.GetUserByEmail(email)
		if err != nil {
			ctx.JSON(422, map[string]interface{}{
				"error": err.Error(),
			})
			ctx.Abort()
			return
		}

		ok := alreadyAdded[u.Email]
		if !ok {
			users = append(users, u)
			alreadyAdded[u.Email] = true
		}
	}

	chats, _ := c.DataRepo.CreateConversation(users, name)

	data := map[string]interface{}{
		"type":    "new_room",
		"room_id": chats.Id,
		"chats": chats,
	}

	payloadInJson, _ := json.Marshal(data)
	for _, p := range chats.Participants {
		// only send to other user
		if p.Id != user.Id {
			c.Publisher.Publish(p.Id, 2, false, string(payloadInJson))
		}
	}

	ctx.JSON(200, map[string]interface{}{
		"user": map[string]interface{}{
			"email": user.Email,
		},
		"chats": chats,
	})
}

func (c *ChatHandler) PostMessage(ctx *gin.Context) {
	roomId := ctx.PostForm("room_id")
	message := ctx.PostForm("message")

	if roomId == "" || message == "" {
		ctx.JSON(422, map[string]interface{}{
			"error": "room id and message cannot be empty",
		})
		ctx.Abort()
		return
	}

	user := ctx.MustGet("user").(*model.User)
	room, err := c.DataRepo.GetConversationById(roomId, user)
	if err != nil {
		ctx.JSON(422, map[string]interface{}{
			"error": err.Error(),
		})
		ctx.Abort()
		return
	}

	chatMessage, err := c.DataRepo.PostMessage(room, user, message)
	if err != nil {
		ctx.JSON(422, map[string]interface{}{
			"error": err.Error(),
		})
		ctx.Abort()
		return
	}

	// publish message to another user, topic is their id
	data := map[string]interface{}{
		"type":    "new_message",
		"room_id": room.Id,
		"message": message,
	}

	payloadInJson, _ := json.Marshal(data)
	for _, p := range room.Participants {
		// only send to other user
		if p.Id != user.Id {
			c.Publisher.Publish(p.Id, 2, false, string(payloadInJson))
		}
	}

	ctx.JSON(200, map[string]interface{}{
		"user": user,
		"chat": chatMessage,
	})
}

func (c *ChatHandler) GetMessageInRoomId(ctx *gin.Context) {
	roomId := ctx.Query("room_id")

	if roomId == "" {
		ctx.JSON(422, map[string]interface{}{
			"error": "room id cannot be empty",
		})
		ctx.Abort()
		return
	}

	user := ctx.MustGet("user").(*model.User)
	room, err := c.DataRepo.GetConversationById(roomId, user)
	if err != nil {
		ctx.JSON(422, map[string]interface{}{
			"error": err.Error(),
		})
		ctx.Abort()
		return
	}

	chatMessages, err := c.DataRepo.GetMessages(room)
	if err != nil {
		ctx.JSON(422, map[string]interface{}{
			"error": err.Error(),
		})
		ctx.Abort()
		return
	}

	ctx.JSON(200, map[string]interface{}{
		"user":     user,
		"messages": chatMessages,
	})
}
