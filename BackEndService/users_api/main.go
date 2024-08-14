package main

import (
	"taxarific_users_api/api"
	"taxarific_users_api/data"
	"github.com/gin-gonic/gin"
	middleware "github.com/oapi-codegen/gin-middleware"
)

func newServer(userAPI *api.API) *gin.Engine {
	swagger, err := api.GetSwagger()
	if err != nil {
		panic(err)
	}
	router := gin.Default()
	router.Use(middleware.OapiRequestValidator(swagger))
	api.RegisterHandlers(router, userAPI)
	return router
}

func main() {
	err := data.NewDB()
	if err != nil {
		panic(err)
	}
	server := newServer(api.NewAPI())
	server.Run(":8080")
}
