package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("postgres", "postgres://JavierC@localhost/todoapp?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

func TodoIndex(w http.ResponseWriter, r *http.Request) {

	rows, err := db.Query("SELECT * FROM todos")
	if err != nil {
		log.Printf("Error opening DB")
		http.Error(w, http.StatusText(500), 500)
		return
	}
	defer rows.Close()

	todos := make([]*Todo, 0)
	for rows.Next() {
		todo := new(Todo)
		err := rows.Scan(&todo.Id, &todo.Name, &todo.Completed)
		if err != nil {
			log.Printf("Error scanning DB")
			http.Error(w, http.StatusText(500), 500)
			return
		}
		todos = append(todos, todo)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error closing DB")
		http.Error(w, http.StatusText(500), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(todos); err != nil {
		log.Printf("Error printing todos")
		panic(err)
	}
}

func TodoCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	name := r.FormValue("name")
	completedString := r.FormValue("completed")

	if name == "" || completedString == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	var completed bool

	if completedString == "true" {
		completed = true
	} else {
		completed = false
	}

	id, err := strconv.ParseInt(r.FormValue("id"), 0, 32)
	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	result, err := db.Exec("INSERT INTO todos VALUES($1, $2, $3)", id, name, completed)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	fmt.Fprintf(w, "Todo %d created successfully (%d row affected)\n", id, rowsAffected)
}

func TodoShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	todoId := vars["todoId"]
	fmt.Fprintln(w, "Todo show:", todoId)
}
