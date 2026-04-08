package api

import (
	"testing"

	"sanitary-warehouse-server/models"
)

// Тест с использованием мока
func TestMockGetProducts(t *testing.T) {
	mockDB := NewMockDB()
	products := mockDB.GetProducts()

	if len(products) == 0 {
		t.Error("Expected at least one product in mock DB")
	}

	if products[0].SKU != "MOCK-001" {
		t.Errorf("Expected SKU MOCK-001, got %s", products[0].SKU)
	}
}

func TestMockCreateProduct(t *testing.T) {
	mockDB := NewMockDB()

	newProduct := models.Product{
		SKU:      "MOCK-003",
		Name:     "Mock Product 3",
		Quantity: 15,
	}

	created := mockDB.CreateProduct(newProduct)

	if created.ID == 0 {
		t.Error("Expected created product to have ID")
	}

	if len(mockDB.GetProducts()) != 3 {
		t.Errorf("Expected 3 products, got %d", len(mockDB.GetProducts()))
	}
}
