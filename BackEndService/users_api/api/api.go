package api

import (
	"taxarific_users_api/auth"
	"taxarific_users_api/data"
	"taxarific_users_api/models"
	"taxarific_users_api/mq"

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
	var err error
	// !! TEST FIRST
	// claim, err := auth.ValidateJWTToken(c.GetHeader("Authorization"))
	// if err != nil {
	// 	c.JSON(401, gin.H{"error": err.Error()})
	// 	return
	// }
	// if claim.Role != "admin" {
	// 	c.JSON(401, gin.H{"error": "unauthorized"})
	// 	return
	// }
	var admin models.PostAdminJSONRequestBody
	if err := c.BindJSON(&admin); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	admin.Password, err = auth.HashPassword(admin.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	err = data.CreateAdmin(&admin)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	user := models.User{
		Email: admin.Email,
	}
	mq.AccountCreationEmail(&user)
	c.JSON(201, gin.H{"message": "admin created"})
}

// PostAdminEmployee implements ServerInterface.
func (a *API) PostAdminEmployee(c *gin.Context) {
	var err error
	// !! Add auth
	var employee models.PostAdminEmployeeJSONRequestBody
	if err := c.BindJSON(&employee); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	employee.Password, err = auth.HashPassword(employee.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	err = data.CreateEmployee(&employee)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
}

// PostLogin implements ServerInterface.
func (a *API) PostLogin(c *gin.Context) {
	var err error
	var login models.PostLoginJSONRequestBody
	if err := c.BindJSON(&login); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	login.Password, err = auth.HashPassword(login.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if *login.Role == "user" {
		user, err := data.Userlogin(string(login.Email))
		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			return
		}
		err = auth.CheckPassword(*user.Password, login.Password)
		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			return
		}
		token, err := auth.GenerateJWTToken(user.Id.Hex(), *login.Role)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"token": token})
		return
	}
	if *login.Role == "admin" {
		admin, err := data.AdminLogin(string(login.Email))
		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			return
		}
		err = auth.CheckPassword(*admin.Password, login.Password)
		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			return
		}
		token, err := auth.GenerateJWTToken(admin.Id.Hex(), *login.Role)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"token": token})
		return
	}
	if *login.Role == "employee" {
		employee, err := data.EmployeeLogin(string(login.Email))
		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			return
		}
		err = auth.CheckPassword(*employee.Password, login.Password)
		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			return
		}
		token, err := auth.GenerateJWTToken(employee.Id.Hex(), *login.Role)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"token": token})
		return
	}
	c.JSON(400, gin.H{"error": "role is required"})
}

// PostUser implements ServerInterface.
func (a *API) PostUser(c *gin.Context) {
	var err error
	var user models.PostUserJSONRequestBody
	if err := c.BindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	user.Password, err = auth.HashPassword(user.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	id, err := data.CreateUser(&user)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	token, err := auth.GenerateJWTToken(id, "user")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error})
		return
	}
	c.JSON(201, gin.H{"token": token})
}

// PutEmployeeAddcaseCaseid implements ServerInterface.
func (a *API) PutEmployeeAddcaseCaseid(c *gin.Context, caseid string) {
	claim, err := auth.ValidateJWTToken(c.GetHeader("Authorization"))
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}
	if claim.Role != "employee" {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	employee, err := data.AddCaseToEmployee(claim.UserId, caseid)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "cases added", "cases": employee.Cases})
}

// PutUserUserid implements ServerInterface.
func (a *API) PutUserUserid(c *gin.Context, userid string) {
	panic("unimplemented")
}

func NewAPI() *API {
	return &API{}
}
