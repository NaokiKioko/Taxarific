package main

import (
	"encoding/json"
	"errors"
	"html/template"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var usersBackendURl string = "http://backend-users-api:8080"

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

	req, err := http.NewRequest("POST", usersBackendURl+"/login", body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute request"})
		return
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
		return
	}
	c.SetCookie("jwt", string(resBody), 3600, "/", "localhost", false, true)
	
	response := map[string]string{"token": string(resBody)}
	c.JSON(http.StatusOK, response)
}

func makeUserRequestBody(c *gin.Context) (map[string]string, error) {
	requestBody := make(map[string]string)

	if username := c.PostForm("username"); username != "" {
		requestBody["username"] = username
	} else {
		c.String(http.StatusBadRequest, "username is required")
		return nil, errors.New("username is required")
	}

	if password := c.PostForm("password"); password != "" {
		requestBody["password"] = password
	} else {
		c.String(http.StatusBadRequest, "password is required")
		return nil, errors.New("password is required")
	}

	if email := c.PostForm("email"); email != "" {
		requestBody["email"] = email
	}
	if phone := c.PostForm("phone"); phone != "" {
		requestBody["phone"] = phone
	}
	if address := c.PostForm("address"); address != "" {
		requestBody["address"] = address
	}
	if city := c.PostForm("city"); city != "" {
		requestBody["city"] = city
	}
	if state := c.PostForm("state"); state != "" {
		requestBody["state"] = state
	}
	if zip := c.PostForm("zip"); zip != "" {
		requestBody["zip"] = zip
	}
	return requestBody, nil
}

func handleSignupPost(c *gin.Context) {

	requestBody, err := makeUserRequestBody(c)
	if err != nil {
		return
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to marshal request body")
		return
	}

	req, err := http.NewRequest("POST", usersBackendURl+"/user", strings.NewReader(string(body)))
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to create request")
		return
	}

	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to execute request")
		return
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to read response body")
		return
	}

	response := map[string]string{"token": string(resBody)}
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
