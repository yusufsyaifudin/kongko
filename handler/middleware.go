package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/yusufsyaifudin/kongko/repo"
)

func ProtectedResource(dataStore repo.DataRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessToken := ctx.Request.Header.Get("ACCESS_TOKEN")

		if accessToken == "" {
			ctx.JSON(400, map[string]interface{}{
				"status": 400,
				"error": "access token header must present",
			})
			ctx.Abort()
			return
		}

		user, err := dataStore.ValidateToken(accessToken)
		if err != nil {
			ctx.JSON(400, map[string]interface{}{
				"status": 400,
				"error": err.Error(),
			})
			ctx.Abort()
			return
		}

		ctx.Set("user", user)

		ctx.Next()
	}
}
