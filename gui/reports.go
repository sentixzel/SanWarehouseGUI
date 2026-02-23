package gui

import (
	"fmt"
	"image/color"
	//"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"SanWarehouse/database"
	"SanWarehouse/models"
)

type Reports struct {
	mainWindow *MainWindow
	window     fyne.Window
}

func NewReports(mw *MainWindow) *Reports {
	return &Reports{
		mainWindow: mw,
	}
}

// ShowReportsMenu показывает меню отчетов
func (r *Reports) ShowReportsMenu() {
	r.window = r.mainWindow.app.NewWindow("Отчеты")
	r.window.Resize(fyne.NewSize(800, 600))

	// Заголовок
	title := canvas.NewText("Отчеты и аналитика", color.White)
	title.TextSize = 24
	title.Alignment = fyne.TextAlignCenter

	headerBg := canvas.NewRectangle(&color.NRGBA{R: 70, G: 70, B: 70, A: 255})
	header := container.NewStack(headerBg, container.NewPadded(title))

	// Кнопки отчетов
	reportsList := container.NewVBox(
		widget.NewCard("", "Общий отчет по складу",
			container.NewVBox(
				widget.NewLabel("Полная статистика по всем товарам"),
				widget.NewButtonWithIcon("Сформировать отчет", theme.DocumentIcon(), r.showGeneralReport),
			),
		),
		widget.NewSeparator(),

		widget.NewCard("", "Товары с истекающим сроком",
			container.NewVBox(
				widget.NewLabel("Товары, которые давно не продаются"),
				widget.NewButtonWithIcon("Анализ оборачиваемости", theme.HistoryIcon(), r.showTurnoverReport),
			),
		),
		widget.NewSeparator(),

		widget.NewCard("", "Финансовый отчет",
			container.NewVBox(
				widget.NewLabel("Прибыль, себестоимость, маржинальность"),
				widget.NewButtonWithIcon("Финансовый анализ", theme.ConfirmIcon(), r.showFinancialReport),
			),
		),
		widget.NewSeparator(),

		widget.NewCard("", "Отчет по категориям",
			container.NewVBox(
				widget.NewLabel("Статистика по категориям товаров"),
				widget.NewButtonWithIcon("Анализ категорий", theme.ViewRestoreIcon(), r.showCategoryReport),
			),
		),
		widget.NewSeparator(),

		widget.NewCard("", "Экспорт данных",
			container.NewVBox(
				widget.NewLabel("Выгрузить данные в CSV"),
				widget.NewButtonWithIcon("Экспорт в CSV", theme.DownloadIcon(), r.exportToCSV),
			),
		),
	)

	scroll := container.NewScroll(reportsList)
	scroll.SetMinSize(fyne.NewSize(700, 450))

	// Основной контейнер
	content := container.NewBorder(header, nil, nil, nil, container.NewPadded(scroll))

	r.window.SetContent(content)
	r.window.Show()
}

// showGeneralReport - общий отчет по складу
func (r *Reports) showGeneralReport() {
	var products []models.Product
	database.DB.Find(&products)

	var totalProducts int64
	var totalValue float64
	var totalPurchaseValue float64
	var totalItems int
	var outOfStock int64
	var lowStock int64
	var activeProducts int64

	// Подсчет статистики
	database.DB.Model(&models.Product{}).Count(&totalProducts)
	database.DB.Model(&models.Product{}).Where("is_active = ?", true).Count(&activeProducts)
	database.DB.Model(&models.Product{}).Select("sum(quantity * selling_price)").Scan(&totalValue)
	database.DB.Model(&models.Product{}).Select("sum(quantity * purchase_price)").Scan(&totalPurchaseValue)
	database.DB.Model(&models.Product{}).Select("sum(quantity)").Scan(&totalItems)
	database.DB.Model(&models.Product{}).Where("quantity - reserved_quantity <= 0").Count(&outOfStock)
	database.DB.Model(&models.Product{}).Where("quantity - reserved_quantity < min_stock_level AND quantity - reserved_quantity > 0").Count(&lowStock)

	// Создаем таблицу с товарами
	data := [][]string{
		{"Показатель", "Значение"},
		{"Всего наименований", fmt.Sprintf("%d", totalProducts)},
		{"Активных товаров", fmt.Sprintf("%d", activeProducts)},
		{"Всего единиц товара", fmt.Sprintf("%d", totalItems)},
		{"Общая стоимость продажи", fmt.Sprintf("%.2f руб.", totalValue)},
		{"Общая себестоимость", fmt.Sprintf("%.2f руб.", totalPurchaseValue)},
		{"Потенциальная прибыль", fmt.Sprintf("%.2f руб.", totalValue-totalPurchaseValue)},
		{"Маржинальность", fmt.Sprintf("%.1f%%", (totalValue-totalPurchaseValue)/totalValue*100)},
		{"Товаров в наличии", fmt.Sprintf("%d", totalProducts-outOfStock)},
		{"Нет в наличии", fmt.Sprintf("%d", outOfStock)},
		{"Низкий запас", fmt.Sprintf("%d", lowStock)},
	}

	list := widget.NewTable(
		func() (int, int) {
			return len(data), 2
		},
		func() fyne.CanvasObject {
			return container.NewStack(
				canvas.NewRectangle(color.Transparent),
				widget.NewLabel("Template"),
			)
		},
		func(id widget.TableCellID, obj fyne.CanvasObject) {
			container := obj.(*fyne.Container)
			label := container.Objects[1].(*widget.Label)
			label.SetText(data[id.Row][id.Col])

			// Стилизация заголовков
			if id.Row == 0 {
				label.TextStyle.Bold = true
				bg := canvas.NewRectangle(&color.NRGBA{R: 200, G: 200, B: 200, A: 255})
				container.Objects[0] = bg
			}

			// ИСПРАВЛЕНО: Убрано обращение к label.Color
		})

	list.SetColumnWidth(0, 250)
	list.SetColumnWidth(1, 200)

	dialog.ShowCustom("Общий отчет по складу", "Закрыть", list, r.mainWindow.window)
}

// showFinancialReport - финансовый отчет
func (r *Reports) showFinancialReport() {
	type CategoryFinance struct {
		Category    string
		Items       int
		PurchaseSum float64
		SellingSum  float64
		Profit      float64
	}

	var results []CategoryFinance

	database.DB.Model(&models.Product{}).
		Select("category, sum(quantity) as items, sum(quantity * purchase_price) as purchase_sum, sum(quantity * selling_price) as selling_sum").
		Group("category").
		Scan(&results)

	// Заголовок
	content := container.NewVBox(
		widget.NewLabelWithStyle("Финансовый анализ по категориям", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
	)

	totalPurchase := 0.0
	totalSelling := 0.0
	totalProfit := 0.0

	for _, r := range results {
		profit := r.SellingSum - r.PurchaseSum
		margin := 0.0
		if r.SellingSum > 0 {
			margin = profit / r.SellingSum * 100
		}

		card := widget.NewCard(r.Category, fmt.Sprintf("Единиц: %d", r.Items),
			container.NewVBox(
				widget.NewLabel(fmt.Sprintf("Закупка: %.2f руб.", r.PurchaseSum)),
				widget.NewLabel(fmt.Sprintf("Продажа: %.2f руб.", r.SellingSum)),
				widget.NewLabelWithStyle(fmt.Sprintf("Прибыль: %.2f руб.", profit), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				widget.NewLabel(fmt.Sprintf("Маржа: %.1f%%", margin)),
			),
		)

		content.Add(card)
		content.Add(widget.NewSeparator())

		totalPurchase += r.PurchaseSum
		totalSelling += r.SellingSum
		totalProfit += profit
	}

	// Итоги
	totalMargin := 0.0
	if totalSelling > 0 {
		totalMargin = totalProfit / totalSelling * 100
	}

	summary := widget.NewCard("ИТОГО", "Общие показатели",
		container.NewVBox(
			widget.NewLabel(fmt.Sprintf("Общая себестоимость: %.2f руб.", totalPurchase)),
			widget.NewLabel(fmt.Sprintf("Общая выручка: %.2f руб.", totalSelling)),
			widget.NewLabelWithStyle(fmt.Sprintf("Общая прибыль: %.2f руб.", totalProfit), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabel(fmt.Sprintf("Общая маржинальность: %.1f%%", totalMargin)),
		),
	)

	content.Add(summary)

	scroll := container.NewScroll(content)
	scroll.SetMinSize(fyne.NewSize(600, 500))

	dialog.ShowCustom("Финансовый отчет", "Закрыть", scroll, r.mainWindow.window)
}

// showCategoryReport - отчет по категориям
func (r *Reports) showCategoryReport() {
	type CategoryStat struct {
		Category    string
		Count       int
		TotalItems  int
		AvgPrice    float64
		BrandsCount int
	}

	var stats []CategoryStat

	// Получаем статистику по категориям
	database.DB.Model(&models.Product{}).
		Select("category, count(*) as count, sum(quantity) as total_items, avg(selling_price) as avg_price").
		Group("category").
		Scan(&stats)

	// Создаем гистограмму (упрощенную)
	content := container.NewVBox()

	for _, s := range stats {
		// Считаем количество брендов в категории
		var brands []string
		database.DB.Model(&models.Product{}).Where("category = ?", s.Category).Distinct("brand").Pluck("brand", &brands)

		// Создаем визуализацию
		bg := canvas.NewRectangle(&color.NRGBA{R: 100, G: 150, B: 255, A: 100})

		statText := fmt.Sprintf("%s:\n  Товаров: %d | Единиц: %d | Брендов: %d | Средняя цена: %.2f руб.",
			s.Category, s.Count, s.TotalItems, len(brands), s.AvgPrice)

		label := widget.NewLabel(statText)

		item := container.NewStack(bg, container.NewPadded(label))
		content.Add(item)
		content.Add(widget.NewSeparator())
	}

	scroll := container.NewScroll(content)
	scroll.SetMinSize(fyne.NewSize(600, 500))

	dialog.ShowCustom("Отчет по категориям", "Закрыть", scroll, r.mainWindow.window)
}

// showTurnoverReport - отчет по оборачиваемости
func (r *Reports) showTurnoverReport() {
	var products []models.Product
	database.DB.Where("quantity > 0").Find(&products)

	content := container.NewVBox(
		widget.NewLabelWithStyle("Анализ оборачиваемости товаров", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
	)

	// Анализируем товары, которые давно не обновлялись
	threshold := time.Now().AddDate(0, -1, 0) // 1 месяц назад

	for _, p := range products {
		if p.UpdatedAt.Before(threshold) {
			days := int(time.Since(p.UpdatedAt).Hours() / 24)

			warning := canvas.NewRectangle(&color.NRGBA{R: 255, G: 200, B: 0, A: 100})

			text := fmt.Sprintf("%s (%s)\n  Не обновлялся %d дней\n  В наличии: %d шт.\n  SKU: %s",
				p.Name, p.Category, days, p.Quantity, p.SKU)

			label := widget.NewLabel(text)

			content.Add(container.NewStack(warning, container.NewPadded(label)))
			content.Add(widget.NewSeparator())
		}
	}

	if len(content.Objects) <= 2 {
		content.Add(widget.NewLabel("Товаров, требующих внимания, не найдено"))
	}

	scroll := container.NewScroll(content)
	scroll.SetMinSize(fyne.NewSize(600, 500))

	dialog.ShowCustom("Оборачиваемость товаров", "Закрыть", scroll, r.mainWindow.window)
}

// exportToCSV - экспорт данных в CSV
func (r *Reports) exportToCSV() {
	var products []models.Product
	database.DB.Find(&products)

	// Создаем CSV контент
	csv := "ID,SKU,Название,Категория,Бренд,Количество,Доступно,Цена закупки,Цена продажи,Статус,Расположение\n"

	for _, p := range products {
		csv += fmt.Sprintf("%d,%s,%s,%s,%s,%d,%d,%.2f,%.2f,%s,%s\n",
			p.ID, p.SKU, p.Name, p.Category, p.Brand,
			p.Quantity, p.AvailableQuantity(),
			p.PurchasePrice, p.SellingPrice,
			p.Status, p.Location)
	}

	// Показываем диалог сохранения
	dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil || writer == nil {
			return
		}
		defer writer.Close()

		writer.Write([]byte(csv))

		dialog.ShowInformation("Экспорт завершен",
			"Данные успешно экспортированы в CSV",
			r.mainWindow.window)
	}, r.mainWindow.window)
}
