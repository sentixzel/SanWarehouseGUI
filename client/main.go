package main

import (
	"flag"
	"log"

	"sanitary-warehouse-client/gui"
)

func main() {
	// Параметр командной строки для адреса сервера
	serverPtr := flag.String("server", "http://localhost:8080", "URL сервера")
	flag.Parse()

	// Запускаем GUI клиент
	app := gui.NewMainWindow(*serverPtr)
	log.Printf("Подключение к серверу: %s", *serverPtr)
	app.Run()
}
