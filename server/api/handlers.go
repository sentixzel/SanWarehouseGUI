package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"sanitary-warehouse-server/database"
	"sanitary-warehouse-server/models"
)

type API struct {
	router *mux.Router
}

func NewAPI() *API {
	api := &API{
		router: mux.NewRouter(),
	}
	api.setupRoutes()
	return api
}

func (api *API) setupRoutes() {
	// Products
	api.router.HandleFunc("/api/products", api.getProducts).Methods("GET")
	api.router.HandleFunc("/api/products/{id}", api.getProduct).Methods("GET")
	api.router.HandleFunc("/api/products", api.createProduct).Methods("POST")
	api.router.HandleFunc("/api/products/{id}", api.updateProduct).Methods("PUT")
	api.router.HandleFunc("/api/products/{id}", api.deleteProduct).Methods("DELETE")
	api.router.HandleFunc("/api/products/search", api.searchProducts).Methods("GET")

	// Reports
	api.router.HandleFunc("/api/reports/lowstock", api.lowStockReport).Methods("GET")
	api.router.HandleFunc("/api/reports/statistics", api.statistics).Methods("GET")
}

func (api *API) GetRouter() *mux.Router {
	return api.router
}

// Получение всех товаров
func (api *API) getProducts(w http.ResponseWriter, r *http.Request) {
	var products []models.Product
	database.DB.Find(&products)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// Получение одного товара
func (api *API) getProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var product models.Product
	result := database.DB.First(&product, id)

	if result.Error != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

// Создание товара
func (api *API) createProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := database.DB.Create(&product)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

// Обновление товара
func (api *API) updateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var product models.Product
	result := database.DB.First(&product, id)
	if result.Error != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	var updateData models.Product
	err := json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Обновляем поля
	product.SKU = updateData.SKU
	product.Name = updateData.Name
	product.Category = updateData.Category
	product.Brand = updateData.Brand
	product.Description = updateData.Description
	product.Quantity = updateData.Quantity
	product.ReservedQuantity = updateData.ReservedQuantity
	product.PurchasePrice = updateData.PurchasePrice
	product.SellingPrice = updateData.SellingPrice
	product.MinStockLevel = updateData.MinStockLevel
	product.Location = updateData.Location
	product.Weight = updateData.Weight
	product.Dimensions = updateData.Dimensions
	product.Material = updateData.Material
	product.MarketplaceID = updateData.MarketplaceID
	product.IsActive = updateData.IsActive

	database.DB.Save(&product)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

// Удаление товара
func (api *API) deleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	result := database.DB.Delete(&models.Product{}, id)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Поиск товаров
func (api *API) searchProducts(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	var products []models.Product
	database.DB.Where("sku ILIKE ? OR name ILIKE ? OR category ILIKE ?",
		"%"+query+"%", "%"+query+"%", "%"+query+"%").Find(&products)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// Отчет по товарам с низким запасом
func (api *API) lowStockReport(w http.ResponseWriter, r *http.Request) {
	var products []models.Product
	database.DB.Where("quantity - reserved_quantity < min_stock_level AND quantity > 0").Find(&products)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// Статистика
func (api *API) statistics(w http.ResponseWriter, r *http.Request) {
	var totalProducts int64
	var totalValue float64
	var totalPurchaseValue float64
	var totalItems int
	var outOfStock int64
	var lowStock int64

	database.DB.Model(&models.Product{}).Count(&totalProducts)
	database.DB.Model(&models.Product{}).Select("sum(quantity * selling_price)").Scan(&totalValue)
	database.DB.Model(&models.Product{}).Select("sum(quantity * purchase_price)").Scan(&totalPurchaseValue)
	database.DB.Model(&models.Product{}).Select("sum(quantity)").Scan(&totalItems)
	database.DB.Model(&models.Product{}).Where("quantity - reserved_quantity <= 0").Count(&outOfStock)
	database.DB.Model(&models.Product{}).Where("quantity - reserved_quantity < min_stock_level AND quantity - reserved_quantity > 0").Count(&lowStock)

	stats := map[string]interface{}{
		"total_products":       totalProducts,
		"total_items":          totalItems,
		"total_value":          totalValue,
		"total_purchase_value": totalPurchaseValue,
		"potential_profit":     totalValue - totalPurchaseValue,
		"out_of_stock":         outOfStock,
		"low_stock":            lowStock,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
