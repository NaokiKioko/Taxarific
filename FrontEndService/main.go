package main

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Define routes
	router.GET("/", handleIndex)
	router.GET("/tax", handleTax)
	router.GET("/start", handleStart)
	router.GET("/quiz", handleQuiz)
	router.GET("/login", handleLogin)
	router.GET("/signup", handleSignup)
	router.GET("/profilenav", handleProfileNav)

	// Start server
	router.Run("127.0.0.1:3000")
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
