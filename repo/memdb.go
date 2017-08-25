// Implementation of DataRepository interface
package repo

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/yusufsyaifudin/kongko/model"
	"sort"
	"time"
)

var secretKey = "RANDOM-STRING-HERE"

type ClaimPayload struct {
	Audience  string `json:"aud,omitempty"`
	Id        string `json:"jti,omitempty"`
	Subject   string `json:"sub,omitempty"`
	Issuer    string `json:"iss"`
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
	NotBefore int64  `json:"nbf"`
	Email     string `json:"email"`
}

func (claim ClaimPayload) Valid() error {
	return nil
}

func NewMemDB() (dataRepo *DataRepoMemDB) {

	users := []*model.User{}
	chats := []*model.ChatRoom{}
	messages := model.ChatMessages{}

	dataRepo = &DataRepoMemDB{
		User:    users,
		Chat:    chats,
		Message: messages,
	}
	return
}

// must implement all method in repo package
type DataRepoMemDB struct {
	User    []*model.User
	Chat    []*model.ChatRoom
	Message model.ChatMessages
}

func (repo *DataRepoMemDB) GetUserByEmail(email string) (user *model.User, err error) {

	for _, u := range repo.User {
		if u.Email == email {
			user = u
			break
		}
	}

	if user == nil {
		err = fmt.Errorf("User with email %v not found", email)
		return
	}

	return
}

func (repo *DataRepoMemDB) RegisterUser(email string, password string) (user *model.User, err error) {
	user, _ = repo.GetUserByEmail(email)

	if user != nil {
		err = fmt.Errorf("User with email %v already registered", email)
		return
	}

	user = &model.User{
		Id:           uuid.New().String(),
		Email:        email,
		Password:     password,
		RegisteredAt: time.Now(),
	}

	// store the user
	repo.User = append(repo.User, user)

	return
}

func (repo *DataRepoMemDB) CreateConversation(users []*model.User, name string) (chat *model.ChatRoom, err error) {

	for _, user := range users {
		_, err = repo.GetUserByEmail(user.Email)
		if err != nil {
			break
		}
	}

	chat = &model.ChatRoom{
		Id:           uuid.New().String(),
		Name:         name,
		Participants: users,
		CreatedAt:    time.Now(),
	}

	repo.Chat = append(repo.Chat, chat)

	return
}

func (repo *DataRepoMemDB) GetConversationById(roomId string, user *model.User) (chat *model.ChatRoom, err error) {
	for _, room := range repo.Chat {
		if room.Id == roomId {
			chat = room
			return
		}
	}

	_, err = repo.CheckRoomEligibility(chat, user)
	if err != nil {
		return
	}

	if chat == nil {
		err = fmt.Errorf("room with id %v not found", roomId)
		return
	}

	return
}

func (repo *DataRepoMemDB) GetConversationList(user *model.User) (chats []*model.ChatRoom, err error) {

	for _, chat := range repo.Chat {
		ok := false // to make sure no duplicate data to show
		for _, u := range chat.Participants {
			if u.Email == user.Email {
				ok = true
			}
		}

		if ok {
			chats = append(chats, chat)
		}
	}

	return
}

func (repo *DataRepoMemDB) PostMessage(room *model.ChatRoom, sender *model.User, message string) (chatMessage *model.ChatMessage, err error) {
	chatMessage = &model.ChatMessage{
		Id:       uuid.New().String(),
		Sender:   sender,
		ChatRoom: room,
		Message:  message,
		SendAt:   time.Now(),
	}

	_, err = repo.CheckRoomEligibility(room, sender)
	if err != nil {
		return
	}

	// post only when eligible
	repo.Message = append(repo.Message, chatMessage)
	return
}

func (repo *DataRepoMemDB) GetMessages(room *model.ChatRoom) (messages model.ChatMessages, err error) {
	for _, message := range repo.Message {
		if message.ChatRoom.Id == room.Id {
			messages = append(messages, message)
		}
	}

	sort.Sort(messages)
	return
}

func (repo *DataRepoMemDB) CheckRoomEligibility(room *model.ChatRoom, user *model.User) (isEligible bool, err error) {
	for _, participant := range room.Participants {
		if user.Id == participant.Id {
			isEligible = true
			break
		}
	}

	if !isEligible {
		err = fmt.Errorf("user %v is not eligible for this room", user.Email)
	}

	return
}

func (repo *DataRepoMemDB) CreateUserToken(user *model.User) (jwtToken string, err error) {
	now := time.Now().Unix()
	expiredAt := time.Now().Add(time.Minute * 60 * 24 * 365).Unix() // one year

	token := jwt.New(jwt.SigningMethodHS256)
	token.Header = map[string]interface{}{
		"typ": "JWT",
		"alg": jwt.SigningMethodHS256.Name,
	}

	token.Claims = ClaimPayload{
		Issuer:    "kongko",  // iss claim
		IssuedAt:  now,       // iat claim
		NotBefore: now,       // nbf claim
		ExpiresAt: expiredAt, // exp claim
		Email:     user.Email,
	}

	jwtToken, err = token.SignedString([]byte(secretKey)) // sign with app secret key
	return
}

func (repo *DataRepoMemDB) ValidateToken(token string) (user *model.User, err error) {
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secretKey), nil
	})

	if err != nil {
		return
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok || !jwtToken.Valid {
		err = fmt.Errorf("Token is not valid.")
		return
	}

	user, err = repo.GetUserByEmail(claims["email"].(string))
	return
}
