package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"taxarific_users_api/data"
	"taxarific_users_api/models"
)

type API struct{}

// GetAdmin implements ServerInterface.
func (a *API) GetAdmin(c *gin.Context) {
	admins, err := data.GetAdmins()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, admins)
}

// GetEmployee implements ServerInterface.
func (a *API) GetEmployee(c *gin.Context) {
}

// GetUser implements ServerInterface.
func (a *API) GetUser(c *gin.Context) {
	
}

// PostUser implements ServerInterface.
func (a *API) PostUser(c *gin.Context) {
	user := models.User{}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
}

// PutEmployeeEmployeeid implements ServerInterface.
func (a *API) PutEmployeeEmployeeid(c *gin.Context, employeeid uuid.UUID) {
}

// PutUserUserid implements ServerInterface.
func (a *API) PutUserUserid(c *gin.Context, userid uuid.UUID) {
	panic("unimplemented")
}

// !! ADMIN ONLY ENDPOINTS MUST BE LOGGED IN !!
// PostAdmin implements ServerInterface.
func (a *API) PostAdmin(c *gin.Context) {
	var adminRequest models.PostAdminJSONRequestBody
	if err := c.ShouldBindJSON(&adminRequest); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// if adminRequest.Name == "" || adminRequest.Email == "" || adminRequest.Password == "" {
	// ! Fix return status code to according to the error
	// c.JSON(400, gin.H{"error": "Name, Email, and Password are required"})
	// return
	// }
	result, err := data.CreateAdmin(adminRequest)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"id": result})
}

// PostAdminEmployee implements ServerInterface.
func (a *API) PostAdminEmployee(c *gin.Context) {
	panic("unimplemented")
}

func NewAPI() *API {
	return &API{}
}
