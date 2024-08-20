package api

import (
	"taxarific_users_api/data"
	"taxarific_users_api/models"
	"taxarific_users_api/services"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

type API struct{}

// GetAdmin implements ServerInterface.
func (a *API) GetAdmin(c *gin.Context) {
	admins, err := services.GetAdmins()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	renderTemplate(c.Writer, "Users.html", map[string]interface{}{
		"users": admins,
	})
}

// GetEmployee implements ServerInterface.
func (a *API) GetEmployee(c *gin.Context) {
	employees, err := services.GetEmployees()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	renderTemplate(c.Writer, "Employees.html", map[string]interface{}{
		"employees": employees,
	})
}

// GetUser implements ServerInterface.
func (a *API) GetUser(c *gin.Context) {
	users, err := services.GetUsers()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	renderTemplate(c.Writer, "Users.html", map[string]interface{}{
		"users": users,
	})
	
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
func (a *API) PutEmployeeEmployeeid(c *gin.Context, employeeid string) {
	user, err := services.GetEmployee(employeeid)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, user)
}

// PutUserUserid implements ServerInterface.
func (a *API) PutUserUserid(c *gin.Context, userid string) {
	user := models.User{}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	services.PutUser(userid, user)
	c.JSON(200, gin.H{"status": "success"})
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

}

func NewAPI() *API {
	return &API{}
}

func renderTemplate(w http.ResponseWriter, templateName string, data interface{}) {
	t, err := template.ParseFiles("templates/" + templateName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}