//go:build integration

package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"sanitary-warehouse-server/models"
)

const (
	baseURL = "http://localhost:8080"
)

func TestMain(m *testing.M) {
	// Ждем запуска сервера
	time.Sleep(5 * time.Second)

	// Проверяем, доступен ли сервер
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(baseURL + "/api/products")
	if err != nil {
		fmt.Printf("Warning: Server not available: %v\n", err)
		fmt.Println("Skipping integration tests")
		os.Exit(0)
	}
	defer resp.Body.Close()

	// Запускаем тесты
	code := m.Run()
	os.Exit(code)
}

func TestCreateAndGetProduct(t *testing.T) {
	// Создаем тестовый продукт
	product := models.Product{
		SKU:          "TEST-001",
		Name:         "Тестовый товар",
		Category:     "Тест",
		Brand:        "TestBrand",
		Quantity:     10,
		SellingPrice: 1000,
	}

	body, _ := json.Marshal(product)
	resp, err := http.Post(baseURL+"/api/products", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Skipf("Server not available: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", resp.StatusCode)
	}

	// Получаем список продуктов
	resp, err = http.Get(baseURL + "/api/products")
	if err != nil {
		t.Fatalf("Failed to get products: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var products []models.Product
	if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(products) == 0 {
		t.Error("Expected at least one product")
	}
}

func TestSearchProducts(t *testing.T) {
	resp, err := http.Get(baseURL + "/api/products/search?q=тест")
	if err != nil {
		t.Skipf("Server not available: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var products []models.Product
	json.NewDecoder(resp.Body).Decode(&products)
	t.Logf("Found %d products matching search", len(products))
}
