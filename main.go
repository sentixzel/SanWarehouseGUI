package main

import (
	"log"

	db "SanWarehouse/database"
	gui "SanWarehouse/gui"
	//"SanWarehouse/models"
)

func main() {
	// Инициализируем базу данных
	if err := db.InitDB(); err != nil {
		log.Fatal("Ошибка инициализации БД:", err)
	}
	defer db.CloseDB()

	// Запускаем GUI
	app := gui.NewMainWindow()
	app.Run()
}
