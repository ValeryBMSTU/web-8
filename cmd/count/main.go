package main

// некоторые импорты нужны для проверки
import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

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

// обработчики http-запросов
func (h *Handlers) GetCounter(w http.ResponseWriter, r *http.Request) {
	msg, err := h.dbProvider.SelectCounter()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Счетчик: " + strconv.Itoa(msg)))
}

func (h *Handlers) PostCounter(w http.ResponseWriter, r *http.Request) {
	input := struct {
		Msg int `json:"msg"`
	}{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}

	err = h.dbProvider.UpdateCounter(input.Msg)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Изменили счетчик!"))
}

// методы для работы с базой данных
func (dp *DatabaseProvider) SelectCounter() (int, error) {
	var msg int

	row := dp.db.QueryRow("SELECT number FROM counter WHERE id_number = 1")
	err := row.Scan(&msg)
	if err != nil {
		return -1, err
	}

	return msg, nil
}

func (dp *DatabaseProvider) UpdateCounter(msg int) error {
	_, err := dp.db.Exec("UPDATE counter SET number = number + $1 WHERE id_number = 1", msg)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	// Считываем аргументы командной строки
	address := flag.String("address", "127.0.0.1:8082", "адрес для запуска сервера")
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
	http.HandleFunc("/get", h.GetCounter)
	http.HandleFunc("/post", h.PostCounter)

	// Запускаем веб-сервер на указанном адресе
	err = http.ListenAndServe(*address, nil)
	if err != nil {
		log.Fatal(err)
	}
}
