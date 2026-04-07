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
	api.router.HandleFunc("/api/products", api.getProducts).Methods("GET")
	api.router.HandleFunc("/api/products/{id}", api.getProduct).Methods("GET")
	api.router.HandleFunc("/api/products", api.createProduct).Methods("POST")
	api.router.HandleFunc("/api/products/{id}", api.updateProduct).Methods("PUT")
	api.router.HandleFunc("/api/products/{id}", api.deleteProduct).Methods("DELETE")
	api.router.HandleFunc("/api/products/search", api.searchProducts).Methods("GET")
	api.router.HandleFunc("/api/reports/lowstock", api.lowStockReport).Methods("GET")
	api.router.HandleFunc("/api/reports/statistics", api.statistics).Methods("GET")
}

func (api *API) GetRouter() *mux.Router {
	return api.router
}

func (api *API) getProducts(w http.ResponseWriter, r *http.Request) {
	var products []models.Product
	database.DB.Find(&products)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// ИСПРАВЛЕНО: проверяем ошибку Encode
	if err := json.NewEncoder(w).Encode(products); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (api *API) getProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	var product models.Product
	result := database.DB.First(&product, id)

	if result.Error != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(product); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

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

	if err := json.NewEncoder(w).Encode(product); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (api *API) updateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	var product models.Product
	result := database.DB.First(&product, id)
	if result.Error != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	var updateData models.Product
	err = json.NewDecoder(r.Body).Decode(&updateData)
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
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(product); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (api *API) deleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	result := database.DB.Delete(&models.Product{}, id)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (api *API) searchProducts(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Search query is required", http.StatusBadRequest)
		return
	}

	var products []models.Product
	database.DB.Where("sku ILIKE ? OR name ILIKE ? OR category ILIKE ?",
		"%"+query+"%", "%"+query+"%", "%"+query+"%").Find(&products)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(products); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (api *API) lowStockReport(w http.ResponseWriter, r *http.Request) {
	var products []models.Product
	database.DB.Where("quantity - reserved_quantity < min_stock_level AND quantity > 0").Find(&products)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(products); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

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
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(stats); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
