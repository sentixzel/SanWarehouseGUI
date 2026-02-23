package database

import (
    "log"
    "os"
    "path/filepath"
    
    "SanWarehouse/models"
    
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() error {
    // Создаем папку data если её нет
    dataDir := "data"
    if err := os.MkdirAll(dataDir, os.ModePerm); err != nil {
        return err
    }
    
    dbPath := filepath.Join(dataDir, "warehouse.db")
    
    var err error
    DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Silent),
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
    
    log.Println("Database initialized successfully")
    return nil
}

func seedData() {
    products := []models.Product{
        {
            SKU:             "MIX-001",
            Name:            "Смеситель для раковины Grohe Eurosmart",
            Category:        "Смесители",
            Brand:           "Grohe",
            Description:     "Однорычажный смеситель для раковины, хромированный, с керамическим картриджем",
            Quantity:        15,
            ReservedQuantity: 2,
            PurchasePrice:   4500,
            SellingPrice:    7990,
            MinStockLevel:   5,
            Location:        "A-01-01",
            Weight:          1.2,
            Dimensions:      "15x20x25",
            Material:        "Латунь",
            MarketplaceID:   "WB-12345",
            IsActive:        true,
        },
        {
            SKU:             "TOI-002",
            Name:            "Унитаз-компакт Cersanit New Now",
            Category:        "Унитазы",
            Brand:           "Cersanit",
            Description:     "Унитаз-компакт с косым выпуском, сиденье микролифт",
            Quantity:        8,
            ReservedQuantity: 1,
            PurchasePrice:   6500,
            SellingPrice:    10990,
            MinStockLevel:   3,
            Location:        "B-02-03",
            Weight:          26.5,
            Dimensions:      "70x38x80",
            Material:        "Фарфор",
            MarketplaceID:   "WB-67890",
            IsActive:        true,
        },
        {
            SKU:             "SINK-003",
            Name:            "Раковина накладная Jacob Delafon Patio",
            Category:        "Раковины",
            Brand:           "Jacob Delafon",
            Description:     "Накладная раковина, 60 см, белая, с переливом",
            Quantity:        3,
            ReservedQuantity: 0,
            PurchasePrice:   5200,
            SellingPrice:    8990,
            MinStockLevel:   4,
            Location:        "C-01-02",
            Weight:          9.5,
            Dimensions:      "60x49x16",
            Material:        "Керамика",
            MarketplaceID:   "WB-24680",
            IsActive:        true,
        },
        {
            SKU:             "BID-004",
            Name:            "Биде Laparet Classic",
            Category:        "Биде",
            Brand:           "Laparet",
            Description:     "Напольное биде, белый глянец",
            Quantity:        2,
            ReservedQuantity: 0,
            PurchasePrice:   3800,
            SellingPrice:    6490,
            MinStockLevel:   2,
            Location:        "B-03-01",
            Weight:          18.0,
            Dimensions:      "40x60x45",
            Material:        "Фарфор",
            MarketplaceID:   "WB-13579",
            IsActive:        true,
        },
        {
            SKU:             "ACC-005",
            Name:            "Набор аксессуаров для ванной IDDIS",
            Category:        "Аксессуары",
            Brand:           "IDDIS",
            Description:     "Набор: стакан, мыльница, дозатор",
            Quantity:        25,
            ReservedQuantity: 5,
            PurchasePrice:   1200,
            SellingPrice:    2490,
            MinStockLevel:   10,
            Location:        "D-01-05",
            Weight:          0.8,
            Dimensions:      "30x15x20",
            Material:        "Керамика/стекло",
            MarketplaceID:   "WB-97531",
            IsActive:        true,
        },
    }
    
    for _, p := range products {
        DB.Create(&p)
    }
    
    log.Println("Test data seeded successfully")
}

// CloseDB закрывает соединение с БД
func CloseDB() error {
    sqlDB, err := DB.DB()
    if err != nil {
        return err
    }
    return sqlDB.Close()
}