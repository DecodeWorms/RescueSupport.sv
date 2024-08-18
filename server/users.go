package server

import (
	"net/http"

	"RescueSupport.sv/handlers"
	"RescueSupport.sv/model"
	"github.com/gin-gonic/gin"
)

type User struct {
	user handlers.Users
}

func NewUser(user handlers.Users) User {
	return User{
		user: user,
	}
}

// CreateUser signUp user
func (u User) SignUp() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//parses data from the incoming request body
		var userSignUp model.UserSignUp
		if err := ctx.ShouldBindJSON(userSignUp); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"data": handleServerResponse(401, "failed", "", err, nil)})
			return
		}

		//CreateUser processes creating of user
		if err := u.user.UserSignUp(userSignUp); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"data": handleServerResponse(500, "failed", "", err, nil)})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": handleServerResponse(200, "success", "", nil, nil)})
	}
}

func (u User) Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//Parse data from the incoming request
		var login model.UserLogin
		if err := ctx.ShouldBindJSON(&login); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"data": handleServerResponse(402, "Failed", "", err, nil)})
			return
		}

		//Call login handler to process the user login
		if err := u.user.Login(login); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"data": handleServerResponse(500, "Failed", "", err, nil)})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"data": handleServerResponse(200, "success", "", nil, nil)})
	}
}

func (u User) ChangePassword() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//Parse the incoming request body
		var changPassword model.UserChangePassword
		if err := ctx.ShouldBindJSON(&changPassword); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"data": handleServerResponse(401, "failed", "", err, nil)})
			return
		}

		//Call changePassword handler to process password changes
		if err := u.user.ChangePassword(model.UserChangePassword(changPassword)); err != nil {
			ctx.JSON(http.StatusInternalServerError, handleServerResponse(500, "failed", "", err, nil))
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": handleServerResponse(200, "success", "", nil, nil)})
	}
}

func (u User) UpdatePassword() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//Parse the incoming request body
		var updatePass model.UserUpdatePassword
		if err := ctx.ShouldBindJSON(&updatePass); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"data": handleServerResponse(401, "failed", "", err, nil)})
			return
		}

		//Call UpdatePassword handler to process the request
		if err := u.user.UpdatePassword(updatePass); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"data": handleServerResponse(500, "success", "", nil, nil)})
			return
		}
	}
}
func handleServerResponse(code int, status, token string, error any, object *model.Users) model.ServerResponse {
	return model.ServerResponse{
		Code:   code,
		Status: status,
		Object: object,
		Error:  error,
		Token:  token,
	}
}
