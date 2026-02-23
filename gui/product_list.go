package gui

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"SanWarehouse/database"
	"SanWarehouse/models"
)

type ProductList struct {
	widget.Table
	mainWindow *MainWindow
	products   []models.Product
}

func NewProductList(mw *MainWindow) *ProductList {
	list := &ProductList{
		mainWindow: mw,
		products:   []models.Product{},
	}

	// Загружаем данные
	database.DB.Find(&list.products)

	// ИСПРАВЛЕНО: функция Length должна возвращать (rows, cols)
	list.Length = func() (int, int) {
		return len(list.products), 10 // 10 колонок
	}

	list.CreateCell = func() fyne.CanvasObject {
		// Создаем ячейку с контейнером для фона и текста
		bg := canvas.NewRectangle(color.Transparent)
		text := widget.NewLabel("Template")
		text.Alignment = fyne.TextAlignCenter
		return container.NewStack(bg, text)
	}

	list.UpdateCell = func(id widget.TableCellID, obj fyne.CanvasObject) {
		if id.Row < 0 || id.Row >= len(list.products) {
			return
		}

		container := obj.(*fyne.Container)
		bg := container.Objects[0].(*canvas.Rectangle)
		label := container.Objects[1].(*widget.Label)

		product := list.products[id.Row]

		// Устанавливаем цвет фона в зависимости от статуса
		switch product.Status {
		case models.StatusLowStock:
			bg.FillColor = &color.NRGBA{R: 255, G: 255, B: 0, A: 50} // Желтый
		case models.StatusOutOfStock:
			bg.FillColor = &color.NRGBA{R: 255, G: 0, B: 0, A: 50} // Красный
		default:
			bg.FillColor = color.Transparent
		}
		bg.Refresh()

		// Устанавливаем текст в зависимости от колонки
		switch id.Col {
		case 0:
			label.SetText(fmt.Sprintf("%d", product.ID))
		case 1:
			label.SetText(product.SKU)
		case 2:
			label.SetText(truncate(product.Name, 20))
		case 3:
			label.SetText(product.Category)
		case 4:
			label.SetText(product.Brand)
		case 5:
			label.SetText(fmt.Sprintf("%d", product.Quantity))
		case 6:
			label.SetText(fmt.Sprintf("%d", product.AvailableQuantity()))
		case 7:
			label.SetText(fmt.Sprintf("%.0f", product.SellingPrice))
		case 8:
			label.SetText(string(product.Status))
		case 9:
			label.SetText(product.Location)
		}
	}

	// ИСПРАВЛЕНО: используем SetColumnWidth вместо прямого доступа
	list.SetColumnWidth(0, 50)  // ID
	list.SetColumnWidth(1, 100) // SKU
	list.SetColumnWidth(2, 200) // Name
	list.SetColumnWidth(3, 100) // Category
	list.SetColumnWidth(4, 100) // Brand
	list.SetColumnWidth(5, 70)  // Quantity
	list.SetColumnWidth(6, 70)  // Available
	list.SetColumnWidth(7, 80)  // Price
	list.SetColumnWidth(8, 100) // Status
	list.SetColumnWidth(9, 80)  // Location

	list.ExtendBaseWidget(list)
	return list
}

func (pl *ProductList) RefreshList() {
	database.DB.Find(&pl.products)
	pl.Refresh()
}

func (pl *ProductList) Search(query string) {
	database.DB.Where("sku LIKE ? OR name LIKE ? OR category LIKE ?",
		"%"+query+"%", "%"+query+"%", "%"+query+"%").Find(&pl.products)
	pl.Refresh()
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

// Добавляем метод для обработки заголовков
func (pl *ProductList) CreateHeader() fyne.CanvasObject {
	headers := []string{"ID", "SKU", "Название", "Категория", "Бренд", "Кол-во", "Доступно", "Цена", "Статус", "Расположение"}

	headerContainer := container.NewWithoutLayout()
	for _, h := range headers {
		label := widget.NewLabelWithStyle(h, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
		headerContainer.Add(label)
	}
	return headerContainer
}

// Добавляем метод для кнопок действий
func (pl *ProductList) DoubleTapped(ev *fyne.PointEvent) {
	// Находим строку по координатам
	for i, product := range pl.products {
		// Простая логика: если координаты в пределах строки
		if ev.Position.Y > float32((i-1)*40+40) && ev.Position.Y < float32((i)*40+40) {
			pl.mainWindow.showProductForm(&product)
			break
		}
	}
}
