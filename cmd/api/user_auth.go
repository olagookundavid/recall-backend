package api

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"firebase.google.com/go/v4/messaging"
	"github.com/gin-gonic/gin"

	"recall-app/cmd/dto"
	"recall-app/internal/domain"
	"recall-app/internal/notification"
	"recall-app/internal/repo"
	"recall-app/internal/token"
)

var (
	EmailResetScope = "email-reset"
	layout          = "02-01-2006"
)

func (app *Application) Test(c *gin.Context) {
	// app.background(func() {
	// 	data := map[string]any{
	// 		"name":            "David",
	// 		"expiryDate":      "4 days",
	// 		"activationToken": "token.Plaintext"}

	// 	err := app.Mailer.Send("dolagookun@icloud.com", EmailTemplate, data)
	// 	if err != nil {
	// 		app.Logger.Error(err.Error(), nil)
	// 	}
	// })

	// go func() {
	notification.SendNotification(app.MessagingClient, c, []string{""}, "", messaging.Notification{
		Title: "Test",
		Body:  "Test",
		// ImageURL: "",
	})
	// }()
}

func (app *Application) GetProfileHandler(c *gin.Context) {
	payload, ok := c.Get(authorizationPayloadKey)
	if !ok {
		app.ServerErrorResponse(c, fmt.Errorf("authorization payload not retrieved successful"))
		return
	}
	tokenPayload := payload.(*token.Payload)
	user, err := app.Handlers.Users.GetById(tokenPayload.UserId)
	if err != nil {
		if err == repo.ErrRecordNotFound {
			app.invalidCredentialsResponse(c)
			return
		}
		app.ServerErrorResponse(c, err)
		return
	}

	rsp := dto.ProfileResponse{
		Message: "User profile returned successfully",
		User:    user.NewUserResponse(),
	}
	c.JSON(http.StatusOK, rsp)

}

func (app *Application) UpdateProfileHandler(c *gin.Context) {
	payload, ok := c.Get(authorizationPayloadKey)
	if !ok {
		app.ServerErrorResponse(c, fmt.Errorf("authorization payload not retrieved successful"))
		return
	}
	tokenPayload := payload.(*token.Payload)

	// c.JSON(http.StatusOK, tokenPayload)
	var req dto.ProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		app.badResponse(c, err.Error())
		return
	}

	//Get By userID
	user, err := app.Handlers.Users.GetById(tokenPayload.UserId)
	if err != nil {
		if err == repo.ErrRecordNotFound {
			app.invalidCredentialsResponse(c)
			return
		}
		app.ServerErrorResponse(c, err)
		return
	}

	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Url != nil {
		user.Url = *req.Url
	}
	if req.Country != nil {
		user.Country = *req.Country
	}
	if req.Phone != nil {
		user.Phone = *req.Phone
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.Dob != nil {
		user.Dob, err = time.Parse(layout, *req.Dob)
		if err != nil {
			app.ServerErrorResponse(c, err)
			return
		}
	}

	err = app.Handlers.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrEditConflict):
			app.editConflictResponse(c)
		default:
			app.ServerErrorResponse(c, err)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})

}

func (app *Application) RegisterUserHandler(c *gin.Context) {
	var req dto.RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		app.badResponse(c, err.Error())
		return
	}

	user := &domain.User{
		Name:    req.Name,
		Phone:   req.Phone,
		Email:   req.Email,
		Country: req.Country,
	}

	err := user.Password.Set(req.Password)
	if err != nil {
		app.ServerErrorResponse(c, err)
		return
	}

	err = app.Handlers.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrDuplicateEmail):
			err = fmt.Errorf("a user with this email address: %s already exists", user.Email)
			app.unAuthorizedResponse(c, err.Error())
		default:
			app.ServerErrorResponse(c, err)
		}
		return
	}

	accessToken, accessPayload, err := getTokenDetails(app, user)

	if err != nil {
		app.ServerErrorResponse(c, err)
		return
	}

	rsp := dto.RegisterUserResponse{
		Message:              "User registered successfully",
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
		// RefreshToken:          refreshToken,
		// RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User: user.NewUserResponse(),
	}
	c.JSON(http.StatusOK, rsp)

}

func (app *Application) LoginUser(c *gin.Context) {
	var req dto.LoginUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		app.badResponse(c, err.Error())
		return
	}

	user, err := app.Handlers.Users.GetByEmail(req.Email)
	if err != nil {
		if err == repo.ErrRecordNotFound {
			app.invalidCredentialsResponse(c)
			return
		}
		app.ServerErrorResponse(c, err)
		return
	}

	match, err := user.Password.Matches(req.Password)
	if err != nil {
		app.ServerErrorResponse(c, err)
		return
	}
	if !match {
		app.invalidCredentialsResponse(c)
		return
	}

	accessToken, accessPayload, err := getTokenDetails(app, user)

	if err != nil {
		app.ServerErrorResponse(c, err)
		return
	}

	rsp := dto.LoginUserResponse{
		Message:              "User logged in successfully",
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
		// RefreshToken:          refreshToken,
		// RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
	}
	c.JSON(http.StatusOK, rsp)
}

func (app *Application) InitiateChangeUserPasswordHandler(c *gin.Context) {
	// Parse and validate the user's new password and password reset token.
	var req struct {
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		app.ServerErrorResponse(c, err)
		return
	}

	//Validate Email

	//Create token in token table

	user, err := app.Handlers.Users.GetByEmail(req.Email)
	if err != nil {
		if err == repo.ErrRecordNotFound {
			app.badResponse(c, "Email not found")
			return
		}
		app.ServerErrorResponse(c, err)
		return
	}

	token := domain.Token{
		UserID: user.ID,
		Email:  req.Email,
		Scope:  EmailResetScope,
		Token:  generateNumericToken(),
		Expiry: time.Now().Add(time.Hour * 24),
	}
	err = app.Handlers.Tokens.Insert(&token)
	if err != nil {
		if err == repo.ErrRecordNotFound {
			app.invalidCredentialsResponse(c)
			return
		}
		app.ServerErrorResponse(c, err)
		return
	}

	//send to email service

	app.background(func() {
		data := map[string]any{
			"name":            user.Name,
			"expiryDate":      token.Expiry.Format("Monday, 02 January 2006 at 15:04"),
			"activationToken": token.Token}

		// err := app.Mailer.Send("dolagookun@icloud.com", "reset-token.html", data)
		err := app.Mailer.Send(user.Email, "reset-token.html", data)
		if err != nil {
			app.Logger.Error(err.Error(), nil)
		}
	})

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Token sent to user email: %s", user.Email)})

}

func (app *Application) UpdatePasswordHandler(c *gin.Context) {
	var req dto.PasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		app.ServerErrorResponse(c, err)
		return
	}

	if req.NewPassword != req.ConfirmPassword {
		app.badResponse(c, "password doesn't match")
		return
	}
	tokenPayload, err := getTokenPayloadFromContext(c)
	if err != nil {
		app.ServerErrorResponse(c, err)
		return
	}

	user, err := app.Handlers.Users.GetById(tokenPayload.UserId)
	if err != nil {
		if err == repo.ErrRecordNotFound {
			app.invalidCredentialsResponse(c)
			return
		}
		app.ServerErrorResponse(c, err)
		return
	}

	match, err := user.Password.Matches(req.OldPassword)
	if err != nil {
		app.ServerErrorResponse(c, err)
		return
	}
	if !match {
		app.invalidCredentialsResponse(c)
		return
	}

	err = user.Password.Set(req.NewPassword)
	if err != nil {
		app.ServerErrorResponse(c, err)
		return
	}
	// Save the updated user record in our database, checking for any edit conflicts as // normal.
	err = app.Handlers.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrEditConflict):
			app.editConflictResponse(c)
		default:
			app.ServerErrorResponse(c, err)
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully changed your password"})
}

func (app *Application) ResetPasswordHandler(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Token    string `json:"token"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		app.ServerErrorResponse(c, err)
		return
	}

	//Check token in token table
	_, err := app.Handlers.Tokens.Get(req.Email, req.Token)
	if err != nil {
		if err == repo.ErrRecordNotFound {
			app.badResponse(c, "Invalid Token")
			return
		}
		app.ServerErrorResponse(c, err)
		return
	}

	user, err := app.Handlers.Users.GetByEmail(req.Email)
	if err != nil {
		if err == repo.ErrRecordNotFound {
			app.invalidCredentialsResponse(c)
			return
		}
		app.ServerErrorResponse(c, err)
		return
	}

	err = user.Password.Set(req.Password)
	if err != nil {
		app.ServerErrorResponse(c, err)
		return
	}
	// Save the updated user record in our database, checking for any edit conflicts as // normal.
	err = app.Handlers.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrEditConflict):
			app.editConflictResponse(c)
		default:
			app.ServerErrorResponse(c, err)
		}
		return
	}

	app.background(func() {
		err = app.Handlers.Tokens.DeleteAllForUser(EmailResetScope, user.ID)
		if err != nil {
			app.ServerErrorResponse(c, err)
			return
		}
	})

	c.JSON(http.StatusOK, gin.H{"message": "Successfully changed your password"})
}

func getTokenDetails(app *Application, user *domain.User) (string, *token.Payload, error) {

	accessDuration, err := time.ParseDuration(app.Config.Token.AccessTokenDuration)
	if err != nil {
		accessDuration = time.Hour * 24
	}

	accessToken, accessPayload, err := app.TokenMaker.CreateToken(user.ID, accessDuration)
	if err != nil {
		return "", nil, err
	}

	// refreshDuration, err := time.ParseDuration(app.Config.Token.RefreshTokenDuration)
	// if err != nil {
	// 	refreshDuration = time.Hour * 168
	// }

	// refreshToken, refreshPayload, err := app.TokenMaker.CreateToken(
	// 	user.ID, refreshDuration,
	// )
	// if err != nil {
	// 	return "", nil, err
	// }
	return accessToken, accessPayload, nil
}

func generateNumericToken() string {

	rand.Seed(time.Now().UnixNano())
	token := ""
	for i := 0; i < 4; i++ {
		token += fmt.Sprintf("%d", rand.Intn(10))
	}
	return token
}

func getTokenPayloadFromContext(c *gin.Context) (*token.Payload, error) {
	payload, ok := c.Get(authorizationPayloadKey)
	if !ok {
		return nil, fmt.Errorf("authorization payload not retrieved successful")
	}
	return payload.(*token.Payload), nil
}
