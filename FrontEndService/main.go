package main

import (
	"html/template"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/tax", handleTax)
	http.HandleFunc("/start", handleStart)
	http.HandleFunc("/quiz", handleQuiz)

	http.ListenAndServe(":3000", nil)
	print("Server started on port 3000")
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index.html", nil)
}

func handleTax(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "tax.html", nil)
}

func handleStart(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "start.html", nil)
}

func handleQuiz(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "quiz.html", nil)
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
