package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

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
	omitempty string 			`json:"omitempty,omitempty"`
}
type JWTResponse struct {
	Token string `json:"token"`
}
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func main() {
	router := gin.Default()

	// Define routes
	router.GET("/", handleIndex)
	router.GET("/tax", handleTax)

	router.GET("/start", handleStart)
	router.GET("/quiz", handleQuiz)

	router.GET("/login", handleUserLogin)
	router.GET("/employee/login", handleRmployeeLogin)
	router.GET("/admin/login", handleAdminLogin)

	router.POST("/login", handleLoginPost)
	router.GET("/logout", handleLogout)
	router.POST("/logout", handleLogoutPost)

	router.GET("/signup", handleSignup)
	router.POST("/signup", handleSignupPost)

	router.GET("/profile", handleProfile)
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

func handleTax(c *gin.Context) {
	renderTemplate(c, "tax.html", nil)
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

	// Decode the response body into the response object
	err = json.NewDecoder(resp.Body).Decode(responseObj)
	if err != nil {
		println("Failed to decode response")
		return fmt.Errorf("failed to decode response: %w", err)
	}
	return nil
}

func HashPassword(password string) (string, error) {
	// bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	// return string(bytes), err
	return password, nil
}
