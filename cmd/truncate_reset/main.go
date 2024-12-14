package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "sandbox"
)

func main() {
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

	// Проверяем соединение
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully connected to the database!")

	// Удаление содержимого таблиц и сброс счетчика последовательности
	tables := []string{"users", "hello", "counter"}
	for _, table := range tables {
		query := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY;", table)
		_, err := db.Exec(query)
		if err != nil {
			log.Fatalf("Failed to truncate table %s: %v", table, err)
		}
		fmt.Printf("Table %s truncated and identity reset successfully.\n", table)
	}
}
