package models

import (
	"testing"
	"time"
)

func TestProductAvailableQuantity(t *testing.T) {
	tests := []struct {
		name     string
		quantity int
		reserved int
		expected int
	}{
		{"Full stock", 10, 0, 10},
		{"Partially reserved", 10, 3, 7},
		{"All reserved", 10, 10, 0},
		{"No stock", 0, 0, 0},
		{"Negative reserved", 10, -1, 11},
		{"Zero quantity", 0, 5, -5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Product{
				Quantity:         tt.quantity,
				ReservedQuantity: tt.reserved,
			}
			got := p.AvailableQuantity()
			if got != tt.expected {
				t.Errorf("AvailableQuantity() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestProductUpdateStatus(t *testing.T) {
	tests := []struct {
		name     string
		quantity int
		reserved int
		minLevel int
		expected ProductStatus
	}{
		{"In stock - plenty", 20, 0, 5, StatusInStock},
		{"In stock - above min", 10, 2, 5, StatusInStock},
		{"Low stock - below min", 4, 0, 5, StatusLowStock},
		{"Low stock - equal to min", 5, 0, 5, StatusInStock},
		{"Out of stock - zero", 0, 0, 5, StatusOutOfStock},
		{"Out of stock - all reserved", 5, 5, 5, StatusOutOfStock},
		{"Low stock with reservation", 7, 3, 5, StatusLowStock},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Product{
				Quantity:         tt.quantity,
				ReservedQuantity: tt.reserved,
				MinStockLevel:    tt.minLevel,
			}
			p.UpdateStatus()
			if p.Status != tt.expected {
				t.Errorf("UpdateStatus() = %v, want %v", p.Status, tt.expected)
			}
		})
	}
}

func TestProductBeforeSave(t *testing.T) {
	tests := []struct {
		name     string
		product  Product
		expected ProductStatus
	}{
		{
			name: "Auto-update status on save",
			product: Product{
				Quantity:      3,
				MinStockLevel: 5,
			},
			expected: StatusLowStock,
		},
		{
			name: "Out of stock auto-update",
			product: Product{
				Quantity:      0,
				MinStockLevel: 5,
			},
			expected: StatusOutOfStock,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Сохраняем оригинальный статус
			originalStatus := tt.product.Status

			// Вызываем BeforeSave
			tt.product.BeforeSave(nil)

			if tt.product.Status == originalStatus && tt.expected != originalStatus {
				t.Errorf("BeforeSave() did not update status. Got %v, want %v",
					tt.product.Status, tt.expected)
			}
		})
	}
}

func TestProductCreation(t *testing.T) {
	p := Product{
		SKU:           "TEST-001",
		Name:          "Test Product",
		Category:      "Test Category",
		Brand:         "Test Brand",
		Description:   "Test Description",
		Quantity:      10,
		SellingPrice:  1000.50,
		PurchasePrice: 800.00,
		Location:      "A-01-01",
		Weight:        1.5,
		Dimensions:    "10x20x30",
		Material:      "Plastic",
		IsActive:      true,
	}

	if p.SKU != "TEST-001" {
		t.Errorf("SKU = %v, want TEST-001", p.SKU)
	}
	if p.Name != "Test Product" {
		t.Errorf("Name = %v, want Test Product", p.Name)
	}
	if p.SellingPrice != 1000.50 {
		t.Errorf("SellingPrice = %v, want 1000.50", p.SellingPrice)
	}
}

func TestProductTimestamps(t *testing.T) {
	now := time.Now()
	p := Product{
		CreatedAt: now,
		UpdatedAt: now,
	}

	if !p.CreatedAt.Equal(now) {
		t.Errorf("CreatedAt = %v, want %v", p.CreatedAt, now)
	}
	if !p.UpdatedAt.Equal(now) {
		t.Errorf("UpdatedAt = %v, want %v", p.UpdatedAt, now)
	}
}

// Бенчмарк тесты
func BenchmarkProductAvailableQuantity(b *testing.B) {
	p := Product{Quantity: 100, ReservedQuantity: 25}

	for i := 0; i < b.N; i++ {
		_ = p.AvailableQuantity()
	}
}

func BenchmarkProductUpdateStatus(b *testing.B) {
	p := Product{Quantity: 10, MinStockLevel: 5}

	for i := 0; i < b.N; i++ {
		p.UpdateStatus()
	}
}
