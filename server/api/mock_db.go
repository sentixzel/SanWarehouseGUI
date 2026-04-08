package api

import (
	"sanitary-warehouse-server/models"
)

type MockDB struct {
	products []models.Product
}

func NewMockDB() *MockDB {
	return &MockDB{
		products: []models.Product{
			{ID: 1, SKU: "MOCK-001", Name: "Mock Product 1", Quantity: 10},
			{ID: 2, SKU: "MOCK-002", Name: "Mock Product 2", Quantity: 5},
		},
	}
}

func (m *MockDB) GetProducts() []models.Product {
	return m.products
}

func (m *MockDB) CreateProduct(product models.Product) models.Product {
	product.ID = uint(len(m.products) + 1)
	m.products = append(m.products, product)
	return product
}
