package database

import (
	"fmt"
	"log"
	"os"

	"sanitary-warehouse-server/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() error {
	// Параметры подключения из переменных окружения
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "warehouse")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return err
	}

	// Автомиграция
	err = DB.AutoMigrate(&models.Product{})
	if err != nil {
		return err
	}

	// Проверяем, нужно ли заполнить тестовыми данными
	var count int64
	DB.Model(&models.Product{}).Count(&count)

	if count == 0 {
		seedData()
	}

	log.Println("Database connected successfully")
	return nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func seedData() {
	products := []models.Product{
		{
			SKU:              "MIX-001",
			Name:             "Смеситель для раковины Grohe Eurosmart",
			Category:         "Смесители",
			Brand:            "Grohe",
			Description:      "Однорычажный смеситель для раковины, хромированный",
			Quantity:         15,
			ReservedQuantity: 2,
			PurchasePrice:    4500,
			SellingPrice:     7990,
			MinStockLevel:    5,
			Location:         "A-01-01",
			Weight:           1.2,
			Dimensions:       "15x20x25",
			Material:         "Латунь",
			MarketplaceID:    "WB-12345",
			IsActive:         true,
		},
		{
			SKU:              "TOI-002",
			Name:             "Унитаз-компакт Cersanit New Now",
			Category:         "Унитазы",
			Brand:            "Cersanit",
			Description:      "Унитаз-компакт с косым выпуском",
			Quantity:         8,
			ReservedQuantity: 1,
			PurchasePrice:    6500,
			SellingPrice:     10990,
			MinStockLevel:    3,
			Location:         "B-02-03",
			Weight:           26.5,
			Dimensions:       "70x38x80",
			Material:         "Фарфор",
			MarketplaceID:    "WB-67890",
			IsActive:         true,
		},
	}

	for _, p := range products {
		DB.Create(&p)
	}

	log.Println("Test data seeded successfully")
}
