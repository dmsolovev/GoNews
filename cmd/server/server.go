package main

import (
	"fmt"
	"log"
	"net/http"

	"GoNews/pkg/api"
	"GoNews/pkg/storage"

	// postgres "GoNews/pkg/storage/postgresql"
	mongodb "GoNews/pkg/storage/mongodb"
)

// Сервер go.
type server struct {
	db  storage.Interface
	api *api.API
}

func main() {
	// Создаём объект сервера.
	var srv server

	// Создаём объекты баз данных.

	// БД в памяти
	// db := memdb.New()
	// fmt.Println("In-memory DB initialized")

	// PostgreSQL
	// pgDB, err := postgres.New("postgres://postgres:1@localhost:5432/GoNews?sslmode=disable")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("PostgreSQL DB initialized")

	// MongoDB
	mongoDB, err := mongodb.New("mongodb://localhost:27017")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("MongoDB initialized")

	// Выбираем какую БД использовать
	// Можно раскомментировать нужную строку:
	//srv.db = db // Использовать in-memory БД
	// srv.db = pgDB // Использовать PostgreSQL
	srv.db = mongoDB // Использовать MongoDB

	// Создаём объект API и регистрируем обработчики.
	srv.api = api.New(srv.db)

	// Запускаем веб-сервер
	if err := http.ListenAndServe(":8080", srv.api.Router()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Starting server on :8080")
}
