package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"html/template"
	"io"
	"net/http"
	"strings"
	"time"

	openapi_types "github.com/oapi-codegen/runtime/types"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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
	router.GET("/profilenav", handleProfileNav)

	// Start server
	router.Run(":3000")
	print("Server started on port 3000")
}

func handleIndex(c *gin.Context) {
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
	email := c.PostForm("email")
	password := c.PostForm("password")

	body := strings.NewReader(`{"username":"` + email + `","password":"` + password + `"}`)

	res, err := sendReqstforJWT(c, usersBackendURl+"/login", body)
	if err != nil {
		return
	}

	c.SetCookie("jwt", res.Token, 3600, "/", domainName, false, true)

	// Respond with the JWT token
	response := map[string]string{"token": res.Token}
	c.JSON(http.StatusOK, response)
}

type User struct {
	Address   *string             `json:"address,omitempty"`
	City      *string             `json:"city,omitempty"`
	CreatedAt *time.Time          `json:"created_at,omitempty"`
	Email     openapi_types.Email `json:"email"`
	Id        primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	Name      string              `json:"name"`
	Password  string              `json:"omitempty"`
	Phone     *string             `json:"phone,omitempty"`
	State     *string             `json:"state,omitempty"`
	UpdatedAt *time.Time          `json:"updated_at,omitempty"`
	Zip       *string             `json:"zip,omitempty"`
}

func handleSignupPost(c *gin.Context) {
	id := primitive.NewObjectID()
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
		Id:       id,
		Address:  &address,
		City:     &city,
		Email:    openapi_types.Email(email),
		Name:     name,
		Password: hashedPassword,
		Phone:    &phone,
		State:    &state,
		Zip:      &zip,
	}

	body, err := userToReader(requestBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal request body"})
		return
	}

	res, err := sendReqstforJWT(c, usersBackendURl+"/signup", body)
	if err != nil {
		return
	}
	c.SetCookie("jwt", res.Token, 3600, "/", domainName, false, true)

	// Respond with the JWT token
	response := map[string]string{"token": res.Token}
	c.JSON(http.StatusOK, response)
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

func userToReader(user User) (io.Reader, error) {
	// Convert the User struct to JSON
	jsonData, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	// Return the JSON as an io.Reader
	return bytes.NewReader(jsonData), nil
}

type JWTResponse struct {
	Token string `json:"token"`
}

func sendReqstforJWT(c *gin.Context, url string, body io.Reader) (JWTResponse, error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return JWTResponse{}, errors.New("failed to create request")
	}

	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute request"})
		return JWTResponse{}, errors.New("failed to execute request")
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
		return JWTResponse{}, errors.New("failed to read response body")
	}

	// Parse the response body into JWTResponse struct
	var jwtResponse JWTResponse
	err = json.Unmarshal(resBody, &jwtResponse)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse JSON response"})
		return JWTResponse{}, errors.New("failed to parse JSON response")
	}

	return jwtResponse, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
