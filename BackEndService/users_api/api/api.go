package api

import (
	"taxarific_users_api/data"
	"taxarific_users_api/models"

	"github.com/gin-gonic/gin"
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
	employees, err := data.GetEmployees()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, employees)
}

// GetUser implements ServerInterface.
func (a *API) GetUser(c *gin.Context) {
	users, err := data.GetUsers()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, users)
}

// PostAdmin implements ServerInterface.
func (a *API) PostAdmin(c *gin.Context) {
	panic("unimplemented")
}

// PostAdminEmployee implements ServerInterface.
func (a *API) PostAdminEmployee(c *gin.Context) {
	panic("unimplemented")
}

// PostLogin implements ServerInterface.
func (a *API) PostLogin(c *gin.Context) {
	var login models.PostLoginJSONRequestBody
	if err := c.BindJSON(&login); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if *login.Role == "user" {
		token, err := data.Userlogin(&login)
		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"token": token})
		return
	}
	if *login.Role == "admin" {
		token, err := data.AdminLogin(&login)
		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"token": token})
		return
	}
	if *login.Role == "employee" {
		token, err := data.EmployeeLogin(&login)
		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"token": token})
		return
	}
	c.JSON(400, gin.H{"error": "role is required"})
	return
}

// PostUser implements ServerInterface.
func (a *API) PostUser(c *gin.Context) {
	var user models.PostUserJSONRequestBody
	if err := c.BindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	token, err := data.CreateUser(&user)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{"token": token})
}

// PutEmployeeAddcaseCaseid implements ServerInterface.
func (a *API) PutEmployeeAddcaseCaseid(c *gin.Context, caseid string) {
	panic("unimplemented")
}

// PutUserUserid implements ServerInterface.
func (a *API) PutUserUserid(c *gin.Context, userid string) {
	panic("unimplemented")
}

func NewAPI() *API {
	return &API{}
}
