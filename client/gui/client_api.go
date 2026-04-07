package gui

import (
	"fmt"

	"github.com/go-resty/resty/v2"

	"sanitary-warehouse-client/models"
)

type ClientAPI struct {
	client  *resty.Client
	baseURL string
}

func NewClientAPI(serverURL string) *ClientAPI {
	client := resty.New()

	return &ClientAPI{
		client:  client,
		baseURL: serverURL,
	}
}

// Получение всех товаров
func (api *ClientAPI) GetProducts() ([]models.Product, error) {
	var products []models.Product

	resp, err := api.client.R().
		SetResult(&products).
		Get(api.baseURL + "/api/products")

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("server returned status: %d", resp.StatusCode())
	}

	return products, nil
}

// Получение товара по ID
func (api *ClientAPI) GetProduct(id uint) (*models.Product, error) {
	var product models.Product

	resp, err := api.client.R().
		SetResult(&product).
		Get(fmt.Sprintf("%s/api/products/%d", api.baseURL, id))

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("server returned status: %d", resp.StatusCode())
	}

	return &product, nil
}

// Создание товара
func (api *ClientAPI) CreateProduct(product *models.Product) error {
	resp, err := api.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(product).
		Post(api.baseURL + "/api/products")

	if err != nil {
		return err
	}

	if resp.StatusCode() != 201 {
		return fmt.Errorf("server returned status: %d", resp.StatusCode())
	}

	return nil
}

// Обновление товара
func (api *ClientAPI) UpdateProduct(id uint, product *models.Product) error {
	resp, err := api.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(product).
		Put(fmt.Sprintf("%s/api/products/%d", api.baseURL, id))

	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("server returned status: %d", resp.StatusCode())
	}

	return nil
}

// Удаление товара
func (api *ClientAPI) DeleteProduct(id uint) error {
	resp, err := api.client.R().
		Delete(fmt.Sprintf("%s/api/products/%d", api.baseURL, id))

	if err != nil {
		return err
	}

	if resp.StatusCode() != 204 {
		return fmt.Errorf("server returned status: %d", resp.StatusCode())
	}

	return nil
}

// Поиск товаров
func (api *ClientAPI) SearchProducts(query string) ([]models.Product, error) {
	var products []models.Product

	resp, err := api.client.R().
		SetResult(&products).
		Get(api.baseURL + "/api/products/search?q=" + query)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("server returned status: %d", resp.StatusCode())
	}

	return products, nil
}

// Получение отчета по низкому запасу
func (api *ClientAPI) GetLowStockReport() ([]models.Product, error) {
	var products []models.Product

	resp, err := api.client.R().
		SetResult(&products).
		Get(api.baseURL + "/api/reports/lowstock")

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("server returned status: %d", resp.StatusCode())
	}

	return products, nil
}

// Получение статистики
func (api *ClientAPI) GetStatistics() (map[string]interface{}, error) {
	var stats map[string]interface{}

	resp, err := api.client.R().
		SetResult(&stats).
		Get(api.baseURL + "/api/reports/statistics")

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("server returned status: %d", resp.StatusCode())
	}

	return stats, nil
}
