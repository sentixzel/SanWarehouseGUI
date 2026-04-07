package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"sanitary-warehouse-server/models"
)

func TestGetProducts(t *testing.T) {
	api := NewAPI()

	req, err := http.NewRequest("GET", "/api/products", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.getProducts)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var products []models.Product
	err = json.NewDecoder(rr.Body).Decode(&products)
	if err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}
}

func TestCreateProduct(t *testing.T) {
	api := NewAPI()

	product := models.Product{
		SKU:          "TEST-UNIT-001",
		Name:         "Unit Test Product",
		Category:     "Testing",
		Brand:        "TestBrand",
		Quantity:     5,
		SellingPrice: 500,
	}

	// ИСПРАВЛЕНО: проверяем ошибку Marshal
	body, err := json.Marshal(product)
	if err != nil {
		t.Fatalf("Failed to marshal product: %v", err)
	}

	req, err := http.NewRequest("POST", "/api/products", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.createProduct)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var createdProduct models.Product
	err = json.NewDecoder(rr.Body).Decode(&createdProduct)
	if err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if createdProduct.SKU != product.SKU {
		t.Errorf("Created product SKU = %v, want %v", createdProduct.SKU, product.SKU)
	}
}

func TestCreateProductInvalidData(t *testing.T) {
	api := NewAPI()

	// Отправляем невалидный JSON
	invalidJSON := []byte(`{"sku": "test", "name":`)

	req, err := http.NewRequest("POST", "/api/products", bytes.NewBuffer(invalidJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.createProduct)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestSearchProducts(t *testing.T) {
	api := NewAPI()

	req, err := http.NewRequest("GET", "/api/products/search?q=test", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.searchProducts)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestSearchProductsEmptyQuery(t *testing.T) {
	api := NewAPI()

	req, err := http.NewRequest("GET", "/api/products/search?q=", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.searchProducts)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestLowStockReport(t *testing.T) {
	api := NewAPI()

	req, err := http.NewRequest("GET", "/api/reports/lowstock", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.lowStockReport)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestStatistics(t *testing.T) {
	api := NewAPI()

	req, err := http.NewRequest("GET", "/api/reports/statistics", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.statistics)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var stats map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&stats)
	if err != nil {
		t.Errorf("Failed to decode statistics: %v", err)
	}

	// Проверяем наличие ключей
	expectedKeys := []string{"total_products", "total_items", "total_value"}
	for _, key := range expectedKeys {
		if _, ok := stats[key]; !ok {
			t.Errorf("Statistics missing key: %s", key)
		}
	}
}

// Дополнительный тест для проверки получения продукта по ID
func TestGetProductByID(t *testing.T) {
	api := NewAPI()

	// Сначала создаем продукт
	product := models.Product{
		SKU:          "TEST-GET-001",
		Name:         "Test Get Product",
		Quantity:     10,
		SellingPrice: 1000,
	}

	body, err := json.Marshal(product)
	if err != nil {
		t.Fatalf("Failed to marshal product: %v", err)
	}

	createReq, _ := http.NewRequest("POST", "/api/products", bytes.NewBuffer(body))
	createReq.Header.Set("Content-Type", "application/json")
	createRR := httptest.NewRecorder()
	http.HandlerFunc(api.createProduct).ServeHTTP(createRR, createReq)

	var createdProduct models.Product
	json.NewDecoder(createRR.Body).Decode(&createdProduct)

	// Теперь получаем продукт по ID
	getReq, _ := http.NewRequest("GET", "/api/products/"+string(rune(createdProduct.ID)), nil)
	getRR := httptest.NewRecorder()
	http.HandlerFunc(api.getProduct).ServeHTTP(getRR, getReq)

	if status := getRR.Code; status != http.StatusOK {
		t.Errorf("Get product returned wrong status: got %v want %v", status, http.StatusOK)
	}
}
