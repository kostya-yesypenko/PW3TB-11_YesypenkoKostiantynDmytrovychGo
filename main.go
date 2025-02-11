package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"go3/handlers"
)

// Головна сторінка
func homePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Помилка завантаження сторінки", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func main() {
	// Налаштування маршрутів
	http.HandleFunc("/", homePage)
	http.HandleFunc("/profitCalc", handlers.ProfitHandler)

	// Підключення стилів
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Запуск сервера
	port := ":8080"
	fmt.Println("http://localhost" + port)
	log.Fatal(http.ListenAndServe(port, nil))
}
