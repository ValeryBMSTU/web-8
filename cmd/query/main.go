package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "89858243234"
	dbname   = "lab_8"
)

type Handlers struct {
	dbProvider DatabaseProvider
}

type DatabaseProvider struct {
	db *sql.DB
}

// Обработчики HTTP-запросов
func (h *Handlers) GetQuery(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Не введен параметр!"))
		return
	}

	test, err := h.dbProvider.SelectQuery(name)
	if !test {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Запись не добавлена в БД!"))
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello," + name + "!"))
}

func (h *Handlers) PostQuery(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Не введен параметр!"))
		return
	}

	test, err := h.dbProvider.SelectQuery(name)
	if test && err == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Запись уже добавлена БД!"))
		return
	}

	err = h.dbProvider.InsertQuery(name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Добавили запись!"))
}

// Методы для работы с базой данных
func (dp *DatabaseProvider) SelectQuery(msg string) (bool, error) {
	var rec string

	row := dp.db.QueryRow("SELECT record FROM query WHERE record = ($1)", msg)
	err := row.Scan(&rec)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (dp *DatabaseProvider) InsertQuery(msg string) error {
	_, err := dp.db.Exec("INSERT INTO query (record) VALUES ($1)", msg)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	// Считываем аргументы командной строки
	address := flag.String("address", "127.0.0.1:8083", "адрес для запуска сервера")
	flag.Parse()

	// Формирование строки подключения для postgres
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Создание соединения с сервером postgres
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Создаем провайдер для БД с набором методов
	dp := DatabaseProvider{db: db}
	// Создаем экземпляр структуры с набором обработчиков
	h := Handlers{dbProvider: dp}

	// Регистрируем обработчики
	http.HandleFunc("/get", h.GetQuery)
	http.HandleFunc("/post", h.PostQuery)

	// Запускаем веб-сервер на указанном адресе
	err = http.ListenAndServe(*address, nil)
	if err != nil {
		log.Fatal(err)
	}
}
