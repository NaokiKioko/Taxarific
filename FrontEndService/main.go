package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	openapi_types "github.com/oapi-codegen/runtime/types"

	"github.com/gin-gonic/gin"
	// "golang.org/x/crypto/bcrypt"
)

var usersBackendURl string = "http://backend-users-api:8080"

// var domainName string = "Taxarific.com"
var domainName string = "http://localhost:3000"

type User struct {
	Id       string              `json:"id,omitempty"`
	Address  string              `json:"address,omitempty"`
	City     string              `json:"city,omitempty"`
	Email    openapi_types.Email `json:"email"`
	Name     string              `json:"name"`
	Password string              `json:"password"`
	Phone    string              `json:"phone,omitempty"`
	State    string              `json:"state,omitempty"`
	Zip      string              `json:"zip,omitempty"`
}
type JWTResponse struct {
	Token string `json:"token"`
}
type LoginRequest struct {
	Username string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
type Case struct {
	Case_id           string `json:"case_id,omitempty"`
	Marital_status    string `json:"marital_status"`
	Dependents        int    `json:"dependents"`
	Employment_status string `json:"employment_status"`
	Estimated_income  string `json:"estimated_income"`
	Case_status       string `json:"case_status,omitempty"`
}

func main() {
	router := gin.Default()

	// Define routes
	router.GET("/", handleIndex)
	router.GET("/case", handleCase)
	// router.POST("/case", handleClaimPost)
	router.GET("/case/view", handleCaseView)

	router.GET("/start", handleStart)
	router.GET("/quiz", handleQuiz)
	router.POST("/quiz", handleQuizPost)

	router.GET("/login", handleUserLogin)
	router.GET("/employee/login", handleRmployeeLogin)
	router.GET("/admin/login", handleAdminLogin)

	router.GET("/employee/create", handleEmplyeeCreate)
	router.GET("/employee", handleEmplyee)
	router.POST("/employee", handleEmployeePost)
	router.GET("/user", handleUser)

	router.POST("/login", handleLoginPost)
	router.GET("/logout", handleLogout)
	router.POST("/logout", handleLogoutPost)

	router.GET("/signup", handleSignup)
	router.POST("/signup", handleSignupPost)

	router.GET("/profile", handleProfile)
	// router.GET("/profile/update", handleProfileUpdate)
	// router.PUT("/profile", handleProfilePut)

	// router.GET("/claim", handleClaim)

	// Start server
	router.Run(":3000")
	print("Server started on port 3000")
}

func handleIndex(c *gin.Context) {
	jwt, err := c.Cookie("JWT")
	if err != nil {
		jwt = ""
	}
	if jwt != "" {
		println("User is logged in")
		renderTemplate(c, "index/user.html", nil)
		return
	}
	println("User is not logged in")
	renderTemplate(c, "index/guest.html", nil)
}

func handleCase(c *gin.Context) {
	renderTemplate(c, "caseSubmit.html", nil)
}

func handleCaseView(c *gin.Context) {
	jwt := c.Query("JWT")
	if jwt == "" {
		renderTemplate(c, "login/userlogin.html", nil)
		return
	}
	Cases := []Case{}
	data := struct {
		Cases []Case
	}{
		Cases: Cases,
	}

	SendRequest("GET", nil, usersBackendURl+"case", data, jwt)
	renderTemplate(c, "cases.html", data)
}

func handleStart(c *gin.Context) {
	role, err := c.Cookie("role")
	if err != nil {
		role = ""
	}
	JWT, err := c.Cookie("JWT")
	if err != nil {
		JWT = ""
	}

	if JWT == "" && role != "" {
		if role == "user" {
			renderTemplate(c, "login/userlogin.html", nil)
			return
		}
		if role == "employee" {
			renderTemplate(c, "login/employeelogin.html", nil)
			return
		}
		if role == "admin" {
			renderTemplate(c, "login/adminlogin.html", nil)
			return
		}
	}

	if JWT != "" && role == "user" {
		renderTemplate(c, "start.html", nil)
		return
	}
	if JWT != "" && (role == "employee" || role == "admin") {
		renderTemplate(c, "nothing.html", nil)
		return
	}

	renderTemplate(c, "login/signup.html", nil)
}

func handleQuiz(c *gin.Context) {
	renderTemplate(c, "quiz.html", nil)
}

func handleQuizPost(c *gin.Context) {

	jwt, err := c.Cookie("JWT")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get JWT"})
		return
	}
	if jwt == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Not logged in"})
	}

	marital_status := c.PostForm("marital_status")
	dependents := c.PostForm("dependents")
	// parse dependents to int
	dependentsInt, err := strconv.Atoi(dependents)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse dependents"})
		return
	}
	employment_status := c.PostForm("employment_status")
	estimated_income := c.PostForm("estimated_income")

	Case := Case{
		Marital_status:    marital_status,
		Dependents:        dependentsInt,
		Employment_status: employment_status,
		Estimated_income:  estimated_income,
	}

	println("Info:")
	println("Marital Status: " + Case.Marital_status)
	print("Dependents: ")
	println(Case.Dependents)
	println("Employment_status: " + Case.Employment_status)
	println("Estimated Income: " + Case.Estimated_income)

	SendRequest("PUT", Case, usersBackendURl+"/user/case", nil, jwt)

	c.JSON(http.StatusOK, nil)

}

func handleUserLogin(c *gin.Context) {
	role, err := c.Cookie("role")
	if err != nil {
		role = ""
	}
	if role != "" {
		if role == "user" {
			println("User loggining in")
			renderTemplate(c, "login/userlogin.html", nil)
			return
		}
		if role == "employee" {
			println("Employee loggining in")
			renderTemplate(c, "login/employeelogin.html", nil)
			return
		}
		if role == "admin" {
			println("Admin loggining in")
			renderTemplate(c, "login/adminlogin.html", nil)
			return
		}
	}
	renderTemplate(c, "login/userlogin.html", nil)
}

func handleRmployeeLogin(c *gin.Context) {
	renderTemplate(c, "login/employeelogin.html", nil)
}

func handleAdminLogin(c *gin.Context) {
	renderTemplate(c, "login/adminlogin.html", nil)
}

func handleLogout(c *gin.Context) {
	renderTemplate(c, "login/logout.html", nil)
}

func handleLogoutPost(c *gin.Context) {
	c.SetCookie("JWT", "", 0, "/", domainName, false, true)
	println("Logged out")
	c.JSON(http.StatusOK, nil)
}

func handleLoginPost(c *gin.Context) {
	jwt, err := c.Cookie("JWT")
	if err != nil {
		jwt = ""
	}
	if jwt != "" {
		c.JSON(http.StatusOK, gin.H{"token": jwt})
		return
	}

	LoginRequest := LoginRequest{
		Username: c.PostForm("email"),
		Password: c.PostForm("password"),
		Role:     c.PostForm("role"),
	}
	println("LoginRequest:")
	println(LoginRequest.Username)
	println(LoginRequest.Password)
	println(LoginRequest.Role)

	JWTResponse := JWTResponse{}
	err = SendRequest("POST", LoginRequest, usersBackendURl+"/login", &JWTResponse, "")
	if err != nil {
		fmt.Println("Failed to send request. Login: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request. Login: " + err.Error()})
		return
	}

	println("JWTResponse Token:")
	println(JWTResponse.Token)

	c.SetCookie("JWT", JWTResponse.Token, 3600, "/", domainName, false, true)
	c.SetCookie("role", "user", 3600, "/", domainName, false, true)
	c.JSON(http.StatusOK, nil)
}

func handleSignupPost(c *gin.Context) {
	// id := primitive.NewObjectID()
	email := c.PostForm("email")
	password := c.PostForm("password")
	name := c.PostForm("name")
	phone := c.PostForm("phone")
	address := c.PostForm("address")
	city := c.PostForm("city")
	state := c.PostForm("state")
	zip := c.PostForm("zip")

	hashedPassword, err := HashPassword(password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	requestBody := User{
		Address:  address,
		City:     city,
		Email:    openapi_types.Email(email),
		Name:     name,
		Password: hashedPassword,
		Phone:    phone,
		State:    state,
		Zip:      zip,
	}

	fmt.Println(requestBody)

	JWTResponse := JWTResponse{}
	err = SendRequest("POST", requestBody, usersBackendURl+"/user", &JWTResponse, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request. Signup: " + err.Error()})
		return
	}

	println("JWTResponse Token:")
	println(JWTResponse.Token)

	c.SetCookie("JWT", JWTResponse.Token, 3600, "/", domainName, false, true)
	c.SetCookie("role", "user", 3600, "/", domainName, false, true)
	c.JSON(http.StatusOK, nil)
}

func handleEmplyee(c *gin.Context) {
	JWT := c.Query("JWT")
	if JWT == "" {
		renderTemplate(c, "login/employeelogin.html", nil)
		return
	}
	role := c.Query("role")
	if role != "admin" {
		renderTemplate(c, "nothing.html", nil)
		return
	}

	Employees := []User{}
	err := SendRequest("GET", nil, usersBackendURl+"/employee", &Employees, JWT)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request. Signup: " + err.Error()})
		return
	}
	data := struct {
		Employees []User
	}{
		Employees: Employees,
	}
	renderTemplate(c, "employee.html", data)
}

func handleUser(c *gin.Context) {
	JWT := c.Query("JWT")
	if JWT == "" {
		renderTemplate(c, "login/userlogin.html", nil)
		return
	}
	role := c.Query("role")
	if role != "admin" {
		renderTemplate(c, "nothing.html", nil)
		return
	}

	Users := []User{}
	err := SendRequest("GET", nil, usersBackendURl+"/user", &Users, JWT)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request. Signup: " + err.Error()})
		return
	}
	data := struct {
		Users []User
	}{
		Users: Users,
	}
	renderTemplate(c, "user.html", data)
}

func handleEmplyeeCreate(c *gin.Context) {
	renderTemplate(c, "employeeCreate.html", nil)
}

func handleEmployeePost(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")
	name := c.PostForm("name")

	hashedPassword, err := HashPassword(password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	requestBody := User{
		Email:    openapi_types.Email(email),
		Name:     name,
		Password: hashedPassword,
	}

	fmt.Println(requestBody)

	err = SendRequest("POST", requestBody, usersBackendURl+"/employee", nil, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request. Signup: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, nil)
}

func handleSignup(c *gin.Context) {
	renderTemplate(c, "login/signup.html", nil)
}

func handleProfile(c *gin.Context) {
	jwt, err := c.Cookie("JWT")
	if err != nil {
		renderTemplate(c, "login/userlogin.html", nil)
		return
	}
	user := User{}
	SendRequest("GET", nil, usersBackendURl+"/user/profile", user, jwt)
	fmt.Print(user)
	renderTemplate(c, "profile.html", user)
}

func renderTemplate(c *gin.Context, templateName string, data interface{}) {
	t, err := template.ParseFiles("templates/" + templateName)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	err = t.Execute(c.Writer, data)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
}

func SendRequest(httpverb string, data interface{}, url string, responseObj interface{}, authHeader string) error {
	// Marshal the request data into JSON

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	// Create a new HTTP POST request
	req, err := http.NewRequest(httpverb, url, bytes.NewBuffer(jsonData))
	if err != nil {
		println(err)
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}
	// Perform the HTTP request
	client := &http.Client{}
	println("Sending request to backend")
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request 1: %w", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		println(resp.StatusCode)
		return fmt.Errorf("failed to send request 2: %w", err)
	}

	if responseObj != nil {
		// Decode the response body into the response object
		err = json.NewDecoder(resp.Body).Decode(responseObj)
		if err != nil {
			println("Failed to decode response")
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}
	return nil
}

func HashPassword(password string) (string, error) {
	// bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	// return string(bytes), err
	return password, nil
}
