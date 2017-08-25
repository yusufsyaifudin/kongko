package repo

import "github.com/yusufsyaifudin/kongko/model"

type DataRepository interface {
	GetUserByEmail(email string) (*model.User, error)
	RegisterUser(email string, password string) (*model.User, error)
	CreateConversation(users []*model.User, name string) (chat *model.ChatRoom, err error)
	GetConversationById(roomId string, user *model.User) (chat *model.ChatRoom, err error)
	GetConversationList(user *model.User) ([]*model.ChatRoom, error)
	PostMessage(room *model.ChatRoom, sender *model.User, message string) (*model.ChatMessage, error)
	GetMessages(room *model.ChatRoom) (model.ChatMessages, error)
	CheckRoomEligibility(room *model.ChatRoom, user *model.User) (isEligible bool, err error)
	CreateUserToken(user *model.User) (token string, err error)
	ValidateToken(token string) (user *model.User, err error)
}
