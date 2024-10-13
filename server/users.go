package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"RescueSupport.sv/config"
	"RescueSupport.sv/handlers"
	"RescueSupport.sv/model"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// var Scopes = []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email", "openid"}
var Scopes = []string{"profile", "email"}
var c = config.ImportConfig(config.OSSource{})

type User struct {
	user handlers.Users
}

func NewUser(user handlers.Users) User {
	return User{
		user: user,
	}
}

func handleOauthConfig(clientID, clientSecret, redirectUrl string, scopes []string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectUrl,
		Endpoint:     google.Endpoint,
		Scopes:       scopes,
	}
}

// CreateUser signUp user
func (u User) SignUp() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//parses data from the incoming request body
		var userSignUp model.UserSignUp
		if err := ctx.ShouldBindJSON(&userSignUp); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"data": handleServerResponse(401, "failed", "", err.Error(), nil)})
			return
		}

		//CreateUser processes creating of user
		if err := u.user.UserSignUp(userSignUp); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"data": handleServerResponse(500, "failed", "", err.Error(), nil)})
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
			ctx.JSON(http.StatusBadRequest, gin.H{"data": handleServerResponse(402, "Failed", "", err.Error(), nil)})
			return
		}

		//Call login handler to process the user login
		if err := u.user.Login(login); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"data": handleServerResponse(500, "Failed", "", err.Error(), nil)})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"data": handleServerResponse(200, "success", "", nil, nil)})
	}
}

func (u User) ChangePassword() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//Get user's ID
		id := ctx.Query("id")
		//Parse the incoming request body
		var changPassword model.UserChangePassword
		if err := ctx.ShouldBindJSON(&changPassword); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"data": handleServerResponse(401, "failed", "", err.Error(), nil)})
			return
		}

		//Call changePassword handler to process password changes
		if err := u.user.ChangePassword(id, changPassword); err != nil {
			ctx.JSON(http.StatusInternalServerError, handleServerResponse(500, "failed", "", err.Error(), nil))
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": handleServerResponse(200, "success", "", nil, nil)})
	}
}

func (u User) UpdatePassword() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//Get user's ID
		id := ctx.Query("id")
		//Parse the incoming request body
		var updatePass model.UserUpdatePassword
		if err := ctx.ShouldBindJSON(&updatePass); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"data": handleServerResponse(401, "failed", "", err.Error(), nil)})
			return
		}

		//Call UpdatePassword handler to process the request
		if err := u.user.UpdatePassword(id, updatePass); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"data": handleServerResponse(500, "failed", "", err.Error(), nil)})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": handleServerResponse(200, "success", "", nil, nil)})
	}
}

func (u User) OauthPage() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		html := `<html><body><a href="/user/oauth_login">Google Login</a></body></html>`
		ctx.Writer.Write([]byte(html))
	}
}

func (u User) LoginWithOauth() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		//Load configuration variables for Oauth
		oauthConfig := handleOauthConfig(c.GoogleClientID, c.GoogleClientSecret, "http://localhost:8001/user/oauth_redirect", Scopes)

		url := oauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
		ctx.Redirect(http.StatusTemporaryRedirect, url)
	}
}

func (u User) GoogleRedirect() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		code := ctx.Query("code") // Use c.Query to get query parameters
		if code == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"data": handleServerResponse(401, "Bad request,Code parameter is missing", "", nil, nil)})
			return
		}

		oauthConfig := handleOauthConfig(c.GoogleClientID, c.GoogleClientSecret, "http://localhost:8001/user/oauth_redirect", Scopes)
		token, err := oauthConfig.Exchange(context.Background(), code)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"data": handleServerResponse(401, "Bad request, failed to exchange token", "", err.Error(), nil)})
			return
		}

		// You can now use the token to access user information.
		client := oauthConfig.Client(context.Background(), token)
		resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"data": handleServerResponse(401, "Bad request,failed to get user info", "", err.Error(), nil)})
			return
		}
		defer resp.Body.Close()
		//userInfo, _ := io.ReadAll(resp.Body)

		//Serialize response body into struct User
		var user model.SignUpWithOauth
		if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
			log.Fatalf("unable to perse the response body %v", err)
			return
		}

		//Store the record into the DB..
		if err := u.user.StoreRecordWithOauth(user); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"data": handleServerResponse(500, "internal server error", "", err.Error(), nil)})
			log.Printf("%d : Internal server error , unable to store user's record %v", 500, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"data": handleServerResponse(200, "success", "", nil, nil)})
	}
}

func (u User) CompleteKyc() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ID := ctx.Query("id")
		if ID == "" {
			ctx.JSON(400, gin.H{"data": handleServerResponse(400, "bad request", "", nil, nil)})
			return
		}

		//Get the request body
		var user model.UserKyc
		if err := ctx.ShouldBindJSON(&user); err != nil {
			ctx.JSON(400, gin.H{"data": handleServerResponse(400, "bad request", "", err.Error(), nil)})
			return
		}

		//handle process of updating user's record
		if err := u.user.CompleteKyc(ID, user); err != nil {
			ctx.JSON(500, gin.H{"data": handleServerResponse(400, "internal server error", "", err.Error(), nil)})
			return
		}
		ctx.JSON(200, gin.H{"data": handleServerResponse(200, "success", "", nil, nil)})
	}
}

func (u User) UpdateKyc() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ID := ctx.Query("id")
		if ID == "" {
			ctx.JSON(400, gin.H{"data": handleServerResponse(400, "bad request", "", nil, nil)})
			return
		}

		//Get the request body
		var user model.UserKyc
		if err := ctx.ShouldBindJSON(&user); err != nil {
			ctx.JSON(400, gin.H{"data": handleServerResponse(400, "bad request", "", err.Error(), nil)})
			return
		}

		//handle process of updating user's record
		if err := u.user.CompleteKyc(ID, user); err != nil {
			ctx.JSON(500, gin.H{"data": handleServerResponse(400, "internal server error", "", err.Error(), nil)})
			return
		}
		ctx.JSON(200, gin.H{"data": handleServerResponse(200, "success", "", nil, nil)})
	}

}

func (u User) ViewProfile() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ID := ctx.Query("id")
		if ID == "" {
			ctx.JSON(400, gin.H{"data": handleServerResponse(400, "bad request", "", nil, nil)})
			return
		}

		//Call user handler to handle fetching of user's record
		res, err := u.user.ViewProfile(ID)
		if err != nil {
			ctx.JSON(500, gin.H{"data": handleServerResponse(500, "internal server error", "", err.Error(), nil)})
			return
		}

		add := model.UserAddress{
			Name:    res.Address.Name,
			City:    res.Address.City,
			Country: res.Address.Country,
			Code:    res.Address.Code,
		}

		//Communicate the response to the client
		ctx.JSON(200, gin.H{"data": handleServerResponse(200, "success", "", nil, &model.Users{
			FirstName: res.FirstName,
			LastName:  res.LastName,
			Email:     res.Email,
			Age:       res.Age,
			Gender:    res.Gender,
			Address:   add,
		})})
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
