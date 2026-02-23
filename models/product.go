package models

import (
	"time"
	"gorm.io/gorm"
)

type ProductStatus string

const (
	StatusInStock ProductStatus = "In stock"
	StatusLowStock ProductStatus = "Low stock"
	StatusOutOfStock ProductStatus = "Out of stock"
	StatusOnOrder ProductStatus = "On order"
)

type Product struct {
    ID              uint           `gorm:"primarykey" json:"id"`
    CreatedAt       time.Time      `json:"created_at"`
    UpdatedAt       time.Time      `json:"updated_at"`
    DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

    SKU             string         `gorm:"uniqueIndex;size:50" json:"sku"`
    Name            string         `gorm:"size:200;not null" json:"name"`
    Category        string         `gorm:"size:100" json:"category"`
    Brand           string         `gorm:"size:100" json:"brand"`
    Description     string         `gorm:"type:text" json:"description"`
    
    Quantity        int            `gorm:"not null;default:0" json:"quantity"`
    ReservedQuantity int           `gorm:"default:0" json:"reserved_quantity"`
    
    PurchasePrice   float64        `json:"purchase_price"`
    SellingPrice    float64        `json:"selling_price"`
    MinStockLevel   int            `gorm:"default:5" json:"min_stock_level"`
    
    Location        string         `gorm:"size:50" json:"location"`
    Status          ProductStatus  `gorm:"size:20;default:'В наличии'" json:"status"`
    
    Weight          float64        `json:"weight"`
    Dimensions      string         `gorm:"size:50" json:"dimensions"`
    Material        string         `gorm:"size:100" json:"material"`
    
    MarketplaceID   string         `gorm:"size:100" json:"marketplace_id"`
    IsActive        bool           `gorm:"default:true" json:"is_active"`
}

func (p *Product) AvailableQuantity() int {
    return p.Quantity - p.ReservedQuantity
}

func (p *Product) UpdateStatus() {
    available := p.AvailableQuantity()
    switch {
    case available <= 0:
        p.Status = StatusOutOfStock
    case available < p.MinStockLevel:
        p.Status = StatusLowStock
    default:
        p.Status = StatusInStock
    }
}

func (p *Product) BeforeSave(tx *gorm.DB) error {
    p.UpdateStatus()
    return nil
}

type ProductDisplay struct {
    ID          uint
    SKU         string
    Name        string
    Category    string
    Brand       string
    Quantity    int
    Available   int
    Price       float64
    Status      ProductStatus
    Location    string
}