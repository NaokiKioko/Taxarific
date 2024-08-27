package api

import (
	"taxarific_users_api/auth"
	"taxarific_users_api/data"
	"taxarific_users_api/models"
	"taxarific_users_api/mq"

	"github.com/gin-gonic/gin"
)

type API struct{}

// GetCase implements ServerInterface.
func (a *API) GetCase(c *gin.Context) {
	cases, err := getAllCases()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, cases)
}

// GetCasePending implements ServerInterface.
func (a *API) GetCasePending(c *gin.Context) {
	cases, err := getAllCases()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	for _, cases := range *cases {
		if *cases.CaseStatus == "pending" {
			c.JSON(200, cases)
		}
	}
	c.JSON(200, cases)
}

// GetUserProfile implements ServerInterface.
func (a *API) GetUserProfile(c *gin.Context) {
	claim, err := auth.ValidateJWTToken(c.GetHeader("Authorization"))
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}
	user, err := data.GetUser(claim.UserId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	*user = models.User{
		Id:    user.Id,
		Email: user.Email,
		Name:  user.Name,
	}
	c.JSON(200, user)
}

// PutUserCase implements ServerInterface.
func (a *API) PutUserCase(c *gin.Context) {
	var putUserCaseRequest models.PutUserCaseJSONRequestBody
	if err := c.BindJSON(&putUserCaseRequest); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	claim, err := auth.ValidateJWTToken(c.GetHeader("Authorization"))
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}
	user, err := data.GetUser(claim.UserId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	user.Case = &putUserCaseRequest.Case
	data.UpdateUser(claim.UserId, user)
	c.JSON(201, gin.H{"message": "case added"})
}

// PutUserProfile implements ServerInterface.
func (a *API) PutUserProfile(c *gin.Context) {
	var updatedUser models.PutUserProfileJSONRequestBody
	if err := c.BindJSON(&updatedUser); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	claim, err := auth.ValidateJWTToken(c.GetHeader("Authorization"))
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}
	user, err := data.GetUser(claim.UserId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if updatedUser.Email != nil {
		user.Email = *updatedUser.Email
	}
	if updatedUser.Password != nil {
		hashedPassword, err := auth.HashPassword(*updatedUser.Password)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		user.Password = &hashedPassword
	}
	if updatedUser.Name != nil {
		user.Name = *updatedUser.Name
	}
	data.UpdateUser(claim.UserId, user)
	c.JSON(201, gin.H{"message": "case added"})
}

// GetAdmin implements ServerInterface.
func (a *API) GetAdmin(c *gin.Context) {
	admins, err := data.GetAdmins()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	adminsResponse := make([]models.UserResponse, len(*admins))
	for i, admin := range *admins {
		adminsResponse[i] = models.UserResponse{
			Email: &admin.Email,
			Id:    &admin.Id,
			Name:  &admin.Name,
		}
	}
	c.JSON(200, adminsResponse)
}

// GetEmployee implements ServerInterface.
func (a *API) GetEmployee(c *gin.Context) {
	employees, err := data.GetEmployees()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	employeeResponses := make([]models.EmployeeResponse, len(*employees))
	for i, employee := range *employees {
		employeeResponses[i] = models.EmployeeResponse{
			Email: &employee.Email,
			Id:    &employee.Id,
			Name:  &employee.Name,
			Cases: employee.Cases,
		}
	}
	c.JSON(200, employeeResponses)
}

// GetUser implements ServerInterface.
func (a *API) GetUser(c *gin.Context) {
	users, err := data.GetUsers()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	userResponses := make([]models.UserResponse, len(*users))
	for i, user := range *users {
		userResponses[i] = models.UserResponse{
			Email: &user.Email,
			Id:    &user.Id,
			Name:  &user.Name,
		}
	}
	c.JSON(200, userResponses)
}

// PostAdmin implements ServerInterface.
func (a *API) PostAdmin(c *gin.Context) {
	claim, err := auth.ValidateJWTToken(c.GetHeader("Authorization"))
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}
	if claim.Role != "admin" {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
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
	claim, err := auth.ValidateJWTToken(c.GetHeader("Authorization"))
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}
	if claim.Role != "admin" {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
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
	c.JSON(201, gin.H{"message": "employee created"})
}

// PostLogin implements ServerInterface.
func (a *API) PostLogin(c *gin.Context) {
	var login models.PostLoginJSONRequestBody
	if err := c.BindJSON(&login); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
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
	cases, err := getAllCases()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	var selectedCase *models.Case
	for i, c := range *cases {
		if *c.CaseId == caseid {
			selectedCase = &(*cases)[i]
		}
	}
	employee, err := data.AddCaseToEmployee(selectedCase, claim.UserId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "cases added", "cases": employee.Cases})
}

func NewAPI() *API {
	return &API{}
}

// Helper methods
func getAllCases() (*[]models.Case, error) {
	users, err := data.GetUsers()
	if err != nil {
		return nil, err
	}
	var cases []models.Case
	for _, user := range *users {
		if user.Case != nil {
			cases = append(cases, *user.Case)
		}
	}
	return &cases, nil
}
