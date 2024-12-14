package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "sandbox"
)

type DatabaseProvider struct {
	db *sql.DB
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	// Подключение к PostgreSQL
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Проверяем соединение
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully connected to the database!")

	// Инициализируем провайдер БД
	dbProvider := &DatabaseProvider{db: db}

	// Регистрируем маршруты
	http.HandleFunc("/api/user", func(w http.ResponseWriter, r *http.Request) {
		userHandler(w, r, dbProvider)
	})

	// Запускаем сервер
	err = http.ListenAndServe(":9000", nil)
	if err != nil {
		log.Fatal(err)
	}
}

// Обработчик запросов
func userHandler(w http.ResponseWriter, r *http.Request, dbProvider *DatabaseProvider) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	switch r.Method {
	case http.MethodGet:
		// Получение информации о пользователе по имени
		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(w, "Parameter 'name' is required", http.StatusBadRequest)
			return
		}

		user, err := dbProvider.GetUser(name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if user == nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		if err := json.NewEncoder(w).Encode(user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	case http.MethodPost:
		// Добавление нового пользователя
		var user User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}

		err = dbProvider.AddUser(user.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "User %s added successfully", user.Name)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Методы работы с базой данных
func (dp *DatabaseProvider) GetUser(name string) (*User, error) {
	query := "SELECT id, name FROM users WHERE name = $1"
	row := dp.db.QueryRow(query, name)

	var user User
	err := row.Scan(&user.ID, &user.Name)
	if err == sql.ErrNoRows {
		return nil, nil // Пользователь не найден
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

func (dp *DatabaseProvider) AddUser(name string) error {
	query := "INSERT INTO users (name) VALUES ($1)"
	_, err := dp.db.Exec(query, name)
	return err
}
