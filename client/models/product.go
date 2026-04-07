package models

import (
	"time"
)

type ProductStatus string

const (
	StatusInStock    ProductStatus = "В наличии"
	StatusLowStock   ProductStatus = "Мало"
	StatusOutOfStock ProductStatus = "Нет в наличии"
	StatusOnOrder    ProductStatus = "Под заказ"
)

type Product struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	SKU         string `json:"sku"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Brand       string `json:"brand"`
	Description string `json:"description"`

	Quantity         int `json:"quantity"`
	ReservedQuantity int `json:"reserved_quantity"`

	PurchasePrice float64 `json:"purchase_price"`
	SellingPrice  float64 `json:"selling_price"`
	MinStockLevel int     `json:"min_stock_level"`

	Location string        `json:"location"`
	Status   ProductStatus `json:"status"`

	Weight     float64 `json:"weight"`
	Dimensions string  `json:"dimensions"`
	Material   string  `json:"material"`

	MarketplaceID string `json:"marketplace_id"`
	IsActive      bool   `json:"is_active"`
}

// AvailableQuantity возвращает доступное количество
func (p *Product) AvailableQuantity() int {
	return p.Quantity - p.ReservedQuantity
}
