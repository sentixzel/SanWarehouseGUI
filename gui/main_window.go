package gui

import (
	"fmt"
	"image/color"

	//"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"SanWarehouse/database"
	"SanWarehouse/models"
)

type MainWindow struct {
	app         fyne.App
	window      fyne.Window
	productList *ProductList
	statusBar   *widget.Label
}

func NewMainWindow() *MainWindow {
	a := app.New()
	w := a.NewWindow("Склад сантехнической гарнитуры")
	w.Resize(fyne.NewSize(1200, 700))

	mw := &MainWindow{
		app:       a,
		window:    w,
		statusBar: widget.NewLabel("Готов к работе"),
	}

	mw.setupUI()
	return mw
}

func (mw *MainWindow) setupUI() {
	// Заголовок
	title := canvas.NewText("Управление складом сантехники", color.White)
	title.TextSize = 20
	title.Alignment = fyne.TextAlignCenter

	headerBg := canvas.NewRectangle(&color.NRGBA{R: 40, G: 40, B: 40, A: 255})
	header := container.NewStack(headerBg, title)

	// Создаем список товаров
	mw.productList = NewProductList(mw)

	// Создаем заголовки таблицы
	headers := mw.productList.CreateHeader()

	// Оборачиваем таблицу в контейнер с заголовками
	tableContainer := container.NewBorder(
		headers, // верхняя часть - заголовки
		nil,
		nil,
		nil,
		container.NewScroll(mw.productList), // центральная часть - таблица с прокруткой
	)

	// Панель инструментов
	toolbar := mw.createToolbar()

	// Основной контент
	content := container.NewBorder(
		container.NewVBox(header, toolbar),
		mw.statusBar,
		nil,
		nil,
		tableContainer,
	)

	mw.window.SetContent(content)

	// Загружаем данные
	mw.productList.RefreshList()
}

func (mw *MainWindow) createToolbar() *widget.Toolbar {
	return &widget.Toolbar{
		Items: []widget.ToolbarItem{
			widget.NewToolbarAction(theme.ContentAddIcon(), func() {
				mw.showProductForm(nil)
			}),
			widget.NewToolbarSeparator(),
			widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {
				mw.productList.RefreshList()
			}),
			widget.NewToolbarSeparator(),
			widget.NewToolbarAction(theme.ContentRedoIcon(), func() {
				mw.showLowStockReport()
			}),
			widget.NewToolbarSeparator(),
			widget.NewToolbarAction(theme.InfoIcon(), func() {
				mw.showStatistics()
			}),
			widget.NewToolbarSeparator(),
			widget.NewToolbarAction(theme.DocumentIcon(), func() {
				reports := NewReports(mw)
				reports.ShowReportsMenu()
			}),
			widget.NewToolbarSeparator(),
			widget.NewToolbarAction(theme.SearchIcon(), func() {
				mw.showSearchDialog()
			}),
		},
	}
}

func (mw *MainWindow) showProductForm(product *models.Product) {
	form := NewProductForm(mw.window, product, func(updatedProduct *models.Product) {
		if product == nil {
			// Создание нового продукта
			database.DB.Create(updatedProduct)
		} else {
			// Обновление существующего
			database.DB.Save(updatedProduct)
		}
		mw.productList.RefreshList()
		mw.statusBar.SetText("Товар сохранен: " + updatedProduct.Name)
	})
	form.Show()
}

func (mw *MainWindow) showLowStockReport() {
	var products []models.Product
	database.DB.Where("quantity - reserved_quantity < min_stock_level AND quantity > 0").Find(&products)

	if len(products) == 0 {
		dialog.ShowInformation("Отчет", "Товаров с низким запасом не найдено", mw.window)
		return
	}

	content := container.NewVBox()
	for _, p := range products {
		available := p.AvailableQuantity()
		text := fmt.Sprintf("%s - %s | Доступно: %d | Мин. уровень: %d",
			p.SKU, p.Name, available, p.MinStockLevel)

		// Создаем цветной индикатор
		bg := canvas.NewRectangle(&color.NRGBA{R: 255, G: 200, B: 200, A: 255})
		label := widget.NewLabel(text)

		content.Add(container.NewStack(bg, container.NewPadded(label)))
	}

	scroll := container.NewScroll(content)
	scroll.SetMinSize(fyne.NewSize(500, 400))

	dialog.ShowCustom("Товары с низким запасом", "Закрыть", scroll, mw.window)
}

func (mw *MainWindow) showStatistics() {
	var totalProducts int64
	var totalValue float64
	var totalItems int
	var outOfStock int64

	database.DB.Model(&models.Product{}).Count(&totalProducts)
	database.DB.Model(&models.Product{}).Select("sum(quantity * purchase_price)").Scan(&totalValue)
	database.DB.Model(&models.Product{}).Select("sum(quantity)").Scan(&totalItems)
	database.DB.Model(&models.Product{}).Where("quantity - reserved_quantity <= 0").Count(&outOfStock)

	stats := fmt.Sprintf(`Статистика склада:
    
    Всего наименований: %d
    Всего единиц товара: %d
    Общая стоимость запасов: %.2f руб.
    Товаров в наличии: %d
    Товаров с нулевым запасом: %d
    Товаров с низким запасом: %d`,
		totalProducts, totalItems, totalValue,
		totalProducts-outOfStock, outOfStock,
		mw.getLowStockCount())

	dialog.ShowInformation("Статистика", stats, mw.window)
}

func (mw *MainWindow) getLowStockCount() int64 {
	var count int64
	database.DB.Model(&models.Product{}).
		Where("quantity - reserved_quantity < min_stock_level AND quantity > 0").
		Count(&count)
	return count
}

func (mw *MainWindow) showSearchDialog() {
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Введите SKU или название товара...")

	items := []*widget.FormItem{
		widget.NewFormItem("Поиск", searchEntry),
	}

	dialog.ShowForm("Поиск товара", "Найти", "Отмена", items, func(b bool) {
		if b {
			query := searchEntry.Text
			mw.productList.Search(query)
			mw.statusBar.SetText("Поиск: " + query)
		}
	}, mw.window)
}

func (mw *MainWindow) Run() {
	mw.window.ShowAndRun()
}
