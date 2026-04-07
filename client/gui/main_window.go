package gui

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"sanitary-warehouse-client/models"
)

type MainWindow struct {
	app         fyne.App
	window      fyne.Window
	productList *ProductList
	statusBar   *widget.Label
	api         *ClientAPI
}

func NewMainWindow(serverURL string) *MainWindow {
	a := app.New()
	w := a.NewWindow("Склад сантехнической гарнитуры (клиент)")
	w.Resize(fyne.NewSize(1200, 700))

	api := NewClientAPI(serverURL)

	mw := &MainWindow{
		app:       a,
		window:    w,
		statusBar: widget.NewLabel("Подключение к серверу..."),
		api:       api,
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

	// Создаем контейнер для таблицы
	tableContainer := container.NewBorder(
		headers,
		nil,
		nil,
		nil,
		container.NewScroll(mw.productList),
	)

	// Панель инструментов
	toolbar := mw.createToolbar()

	// Статус бар с информацией о подключении
	statusContainer := container.NewHBox(
		widget.NewIcon(theme.InfoIcon()),
		mw.statusBar,
		widget.NewLabel(" | "),
		widget.NewLabel("Сервер: "+mw.api.baseURL),
	)

	// Основной контент
	content := container.NewBorder(
		container.NewVBox(header, toolbar),
		statusContainer,
		nil,
		nil,
		tableContainer,
	)

	mw.window.SetContent(content)

	// Загружаем данные
	mw.RefreshData()
}

func (mw *MainWindow) createToolbar() *widget.Toolbar {
	return &widget.Toolbar{
		Items: []widget.ToolbarItem{
			widget.NewToolbarAction(theme.ContentAddIcon(), func() {
				mw.showProductForm(nil)
			}),
			widget.NewToolbarSeparator(),
			widget.NewToolbarAction(theme.ContentRedoIcon(), func() {
				mw.RefreshData()
			}),
			widget.NewToolbarSeparator(),
			widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {
				mw.showLowStockReport()
			}),
			widget.NewToolbarSeparator(),
			widget.NewToolbarAction(theme.InfoIcon(), func() {
				mw.showStatistics()
			}),
			widget.NewToolbarAction(theme.DocumentIcon(), func() {
				// Отчеты будут позже
				dialog.ShowInformation("Информация", "Функция отчетов в разработке", mw.window)
			}),
			widget.NewToolbarSpacer(),
			widget.NewToolbarAction(theme.SearchIcon(), func() {
				mw.showSearchDialog()
			}),
		},
	}
}

func (mw *MainWindow) RefreshData() {
	mw.statusBar.SetText("Загрузка данных...")

	products, err := mw.api.GetProducts()
	if err != nil {
		mw.statusBar.SetText("Ошибка подключения к серверу")
		dialog.ShowError(fmt.Errorf("Не удалось подключиться к серверу: %v", err), mw.window)
		return
	}

	mw.productList.SetProducts(products)
	mw.statusBar.SetText(fmt.Sprintf("Загружено %d товаров", len(products)))
}

func (mw *MainWindow) showProductForm(product *models.Product) {
	form := NewProductForm(mw.window, product, func(updatedProduct *models.Product) {
		var err error
		if product == nil {
			// Создание нового продукта
			err = mw.api.CreateProduct(updatedProduct)
		} else {
			// Обновление существующего
			err = mw.api.UpdateProduct(product.ID, updatedProduct)
		}

		if err != nil {
			dialog.ShowError(fmt.Errorf("Ошибка сохранения: %v", err), mw.window)
			return
		}

		mw.RefreshData()
		mw.statusBar.SetText("Товар сохранен: " + updatedProduct.Name)
	})

	form.Show()
}

func (mw *MainWindow) showLowStockReport() {
	products, err := mw.api.GetLowStockReport()
	if err != nil {
		dialog.ShowError(fmt.Errorf("Ошибка получения отчета: %v", err), mw.window)
		return
	}

	if len(products) == 0 {
		dialog.ShowInformation("Отчет", "Товаров с низким запасом не найдено", mw.window)
		return
	}

	content := container.NewVBox()
	for _, p := range products {
		available := p.AvailableQuantity()
		text := fmt.Sprintf("%s - %s | Доступно: %d | Мин. уровень: %d",
			p.SKU, p.Name, available, p.MinStockLevel)

		bg := canvas.NewRectangle(&color.NRGBA{R: 255, G: 200, B: 200, A: 255})
		label := widget.NewLabel(text)

		content.Add(container.NewStack(bg, container.NewPadded(label)))
	}

	scroll := container.NewScroll(content)
	scroll.SetMinSize(fyne.NewSize(500, 400))

	dialog.ShowCustom("Товары с низким запасом", "Закрыть", scroll, mw.window)
}

func (mw *MainWindow) showStatistics() {
	stats, err := mw.api.GetStatistics()
	if err != nil {
		dialog.ShowError(fmt.Errorf("Ошибка получения статистики: %v", err), mw.window)
		return
	}

	statText := fmt.Sprintf(`Статистика склада:
    
    Всего наименований: %.0f
    Всего единиц товара: %.0f
    Общая стоимость продажи: %.2f руб.
    Общая себестоимость: %.2f руб.
    Потенциальная прибыль: %.2f руб.
    Товаров с нулевым запасом: %.0f
    Товаров с низким запасом: %.0f`,
		stats["total_products"].(float64),
		stats["total_items"].(float64),
		stats["total_value"].(float64),
		stats["total_purchase_value"].(float64),
		stats["potential_profit"].(float64),
		stats["out_of_stock"].(float64),
		stats["low_stock"].(float64))

	dialog.ShowInformation("Статистика", statText, mw.window)
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
			products, err := mw.api.SearchProducts(query)
			if err != nil {
				dialog.ShowError(fmt.Errorf("Ошибка поиска: %v", err), mw.window)
				return
			}
			mw.productList.SetProducts(products)
			mw.statusBar.SetText(fmt.Sprintf("Найдено %d товаров по запросу '%s'", len(products), query))
		}
	}, mw.window)
}

func (mw *MainWindow) Run() {
	mw.window.ShowAndRun()
}
