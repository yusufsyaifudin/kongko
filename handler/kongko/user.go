package kongko

import (
	"github.com/gin-gonic/gin"
	"github.com/yusufsyaifudin/kongko/repo"
)

type UserHandler struct {
	DataRepo repo.DataRepository
}

func NewUserHandler(dataRepo repo.DataRepository) (userHandler *UserHandler) {
	userHandler = &UserHandler{
		DataRepo: dataRepo,
	}
	return
}

func (u *UserHandler) Register(ctx *gin.Context) {
	email := ctx.PostForm("email")
	password := ctx.PostForm("password")

	if email == "" || password == "" {
		ctx.JSON(400, map[string]interface{}{
			"error": "Email and password must not be empty.",
		})
		ctx.Abort()
		return
	}

	user, err := u.DataRepo.RegisterUser(email, password)
	if err != nil {
		ctx.JSON(422, map[string]interface{}{
			"error": err.Error(),
		})
		ctx.Abort()
		return
	}

	ctx.JSON(200, map[string]interface{}{
		"user": map[string]interface{}{
			"id":            user.Id,
			"email":         user.Email,
			"registered_at": user.RegisteredAt,
		},
	})
}

func (u *UserHandler) Login(ctx *gin.Context) {
	email := ctx.PostForm("email")
	password := ctx.PostForm("password")

	if email == "" || password == "" {
		ctx.JSON(400, map[string]interface{}{
			"error": "Email and password must not be empty.",
		})
		ctx.Abort()
		return
	}

	user, err := u.DataRepo.GetUserByEmail(email)
	if err != nil {
		ctx.JSON(422, map[string]interface{}{
			"error": err.Error(),
		})
		ctx.Abort()
		return
	}

	if user.Password != password {
		ctx.JSON(401, map[string]interface{}{
			"error": "Wrong password.",
		})
		ctx.Abort()
		return
	}

	token, err := u.DataRepo.CreateUserToken(user)
	if err != nil {
		ctx.JSON(422, map[string]interface{}{
			"error": err.Error(),
		})
		ctx.Abort()
		return
	}

	ctx.JSON(200, map[string]interface{}{
		"user": map[string]interface{}{
			"id":            user.Id,
			"email":         user.Email,
			"registered_at": user.RegisteredAt,
			"token":         token,
		},
	})
}
