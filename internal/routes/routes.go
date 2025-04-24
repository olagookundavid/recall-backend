package routes

import (
	"expvar"
	"recall-app/cmd/api"

	"github.com/gin-gonic/gin"
)

func Routes(app *api.Application) *gin.Engine {
	r := gin.Default()

	r.Use(gin.CustomRecovery(app.InternalServerErrorHandler))
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
	UserRoutes(v1_api, app)
	ProductRoutes(v1_api, app)

	return r
}

func UserRoutes(r *gin.RouterGroup, app *api.Application) {
	user := r.Group("/user")
	withAuth := user.Group("/", app.TokenMiddleware(app.TokenMaker))

	user.POST("/login", app.LoginUser)
	user.POST("/register", app.RegisterUserHandler)
	user.POST("/getResetToken", app.InitiateChangeUserPasswordHandler)
	user.POST("/resetPassword", app.ResetPasswordHandler)
	withAuth.POST("/profile", app.UpdateProfileHandler)
	withAuth.GET("/profile", app.GetProfileHandler)
	withAuth.POST("/changePassword", app.UpdatePasswordHandler)
	withAuth.POST("/test", app.Test)

}

func ProductRoutes(r *gin.RouterGroup, app *api.Application) {
	product := r.Group("/product")
	withAuth := product.Group("/", app.TokenMiddleware(app.TokenMaker))

	withAuth.GET("/getQrProduct", app.GetProductFromQR)
	//Tracked
	withAuth.POST("/", app.CreateProductHandler)
	withAuth.GET("/", app.GetProductFromQR)
	withAuth.DELETE("/:id", app.DeleteProductHandler)

	//Recalls
	// withAuth.POST("/C", app.CreateProductHandler)
	// withAuth.GET("/product", app.GetProductFromQR)
	// withAuth.DELETE("/product/:id", app.DeleteProductHandler)
}
