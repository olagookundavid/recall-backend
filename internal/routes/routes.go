package routes

import (
	"expvar"
	"recall-app/cmd/api"

	"github.com/gin-gonic/gin"
)

func Routes(app *api.Application) *gin.Engine {
	r := gin.Default()

	r.NoRoute(app.NotFoundResponse)
	r.NoMethod(app.MethodNotAllowedResponse)

	// Root endpoints
	r.GET("/healthcheck", app.HealthcheckHandler)
	r.GET("/debug/vars", app.WrapHTTPHandler(expvar.Handler()))

	// Group: /api/v1/
	v1_api := r.Group("/api/v1", app.RecoverPanic(), app.RateLimit(), app.Metrics())

	// test := r.Group("/")
	// test.GET("/healthcheck", app.HealthcheckHandler)
	// Register subroutes
	userRoutes(v1_api, app)
	secondRoutes(v1_api, app)

	return r
}

func userRoutes(r *gin.RouterGroup, app *api.Application) {
	user := r.Group("/user")
	withAuth := user.Group("/", app.TokenMiddleware(app.TokenMaker))

	user.POST("/login", app.LoginUser)
	user.POST("/register", app.RegisterUserHandler)
	user.POST("/getResetToken", app.InitiateChangeUserPasswordHandler)
	user.POST("/resetPassword", app.ChangeUserPasswordHandler)
	withAuth.POST("/test", app.Test)
	// user.POST("/password-reset", app.CreatePasswordResetTokenHandler)

}

func secondRoutes(r *gin.RouterGroup, app *api.Application) {
	//  := r.Group("/")
	// withAuth := .Group("/", app.TokenMiddleware(app.TokenMaker))

}
