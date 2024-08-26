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

var domainName string = "Taxarific.com"

func main() {
	router := gin.Default()

	// Define routes
	router.GET("/", handleIndex)
	router.GET("/tax", handleTax)
	
	router.GET("/start", handleStart)
	router.GET("/quiz", handleQuiz)

	router.GET("/login", handleLogin)
	router.POST("/login", handleLoginPost)

	router.GET("/signup", handleSignup)
	router.POST("/signup", handleSignupPost)

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
		renderTemplate(c, "index/user.html", nil)
		return
	}
	renderTemplate(c, "index/guest.html", nil)
}

func handleTax(c *gin.Context) {
	renderTemplate(c, "tax.html", nil)
}

func handleStart(c *gin.Context) {
	renderTemplate(c, "start.html", nil)
}

func handleQuiz(c *gin.Context) {
	renderTemplate(c, "quiz.html", nil)
}

func handleLogin(c *gin.Context) {
	renderTemplate(c, "login/login.html", nil)
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
	}
	JWTResponse := JWTResponse{}
	err = SendRequest("GET", LoginRequest, usersBackendURl+"/user", &JWTResponse, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request"})
		return
	}
	c.SetCookie("JWT", "Bearer "+JWTResponse.Token, 3600, "/", domainName, false, true)
	c.JSON(http.StatusOK, gin.H{"token": JWTResponse.Token})
}

type User struct {
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
	Username string `json:"username"`
	Password string `json:"password"`
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
		println("Failed to send request")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request"})
		return
	}

	println("JWTResponse Token:")
	println(JWTResponse.Token)

	c.SetCookie("JWT", "Bearer "+JWTResponse.Token, 3600, "/", domainName, false, true)
	c.JSON(http.StatusOK, gin.H{"token": JWTResponse.Token})
}

func handleSignup(c *gin.Context) {
	renderTemplate(c, "login/signup.html", nil)
}

func handleProfileNav(c *gin.Context) {
	renderTemplate(c, "profilenav.html", nil)
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
		return fmt.Errorf("failed to send request: %w", err)
	}
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send request: %w", err)
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
