package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Data struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	http.HandleFunc("/data", handleData)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleData(w http.ResponseWriter, r *http.Request) {
	// Установка соединения с базой данных
	db, err := sql.Open("mysql", "root:admin@tcp(localhost:3306)/sakila")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Выполнение запроса к базе данных
	rows, err := db.Query("Select actor_id, first_name from sakila.actor")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Сбор данных из запроса
	var data []Data
	for rows.Next() {
		var d Data
		if err := rows.Scan(&d.ID, &d.Name); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data = append(data, d)
	}

	// Конвертирование данных в формат JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Установка заголовка Content-Type и отправка данных в ответе
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
