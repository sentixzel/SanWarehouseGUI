package gui

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"sanitary-warehouse-client/models"
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

	// Функция определения размеров таблицы
	list.Length = func() (int, int) {
		return len(list.products), 10
	}

	// Создание ячейки
	list.CreateCell = func() fyne.CanvasObject {
		bg := canvas.NewRectangle(color.Transparent)

		btnContainer := container.NewHBox(
			widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), nil),
			widget.NewButtonWithIcon("", theme.DeleteIcon(), nil),
		)
		btnContainer.Hide()

		text := widget.NewLabel("Template")
		text.Alignment = fyne.TextAlignCenter

		return container.NewStack(bg, text, btnContainer)
	}

	// Обновление ячейки
	list.UpdateCell = func(id widget.TableCellID, obj fyne.CanvasObject) {
		if id.Row < 0 || id.Row >= len(list.products) {
			return
		}

		stack := obj.(*fyne.Container)
		bg := stack.Objects[0].(*canvas.Rectangle)
		text := stack.Objects[1].(*widget.Label)
		btnContainer := stack.Objects[2].(*fyne.Container)

		product := list.products[id.Row]

		// Цвет фона
		switch product.Status {
		case models.StatusLowStock:
			bg.FillColor = &color.NRGBA{R: 255, G: 255, B: 0, A: 50}
		case models.StatusOutOfStock:
			bg.FillColor = &color.NRGBA{R: 255, G: 0, B: 0, A: 50}
		default:
			bg.FillColor = color.Transparent
		}
		bg.Refresh()

		// Последняя колонка - кнопки
		if id.Col == 9 {
			text.Hide()
			btnContainer.Show()

			editBtn := btnContainer.Objects[0].(*widget.Button)
			editBtn.OnTapped = func() {
				list.mainWindow.showProductForm(&product)
			}

			deleteBtn := btnContainer.Objects[1].(*widget.Button)
			deleteBtn.OnTapped = func() {
				list.showDeleteConfirmation(product)
			}
		} else {
			text.Show()
			btnContainer.Hide()

			switch id.Col {
			case 0:
				text.SetText(fmt.Sprintf("%d", product.ID))
			case 1:
				text.SetText(product.SKU)
			case 2:
				text.SetText(truncate(product.Name, 20))
			case 3:
				text.SetText(product.Category)
			case 4:
				text.SetText(product.Brand)
			case 5:
				text.SetText(fmt.Sprintf("%d", product.Quantity))
			case 6:
				text.SetText(fmt.Sprintf("%d", product.AvailableQuantity()))
			case 7:
				text.SetText(fmt.Sprintf("%.0f", product.SellingPrice))
			case 8:
				text.SetText(string(product.Status))
			}
		}
	}

	// Ширина колонок
	list.SetColumnWidth(0, 50)
	list.SetColumnWidth(1, 100)
	list.SetColumnWidth(2, 200)
	list.SetColumnWidth(3, 100)
	list.SetColumnWidth(4, 100)
	list.SetColumnWidth(5, 70)
	list.SetColumnWidth(6, 70)
	list.SetColumnWidth(7, 80)
	list.SetColumnWidth(8, 100)
	list.SetColumnWidth(9, 100)

	list.ExtendBaseWidget(list)
	return list
}

func (pl *ProductList) SetProducts(products []models.Product) {
	pl.products = products
	pl.Refresh()
}

func (pl *ProductList) CreateHeader() fyne.CanvasObject {
	headers := []string{"ID", "SKU", "Название", "Категория", "Бренд", "Кол-во", "Доступно", "Цена", "Статус", "Действия"}

	headerContainer := container.NewGridWithColumns(10)
	for _, h := range headers {
		label := widget.NewLabelWithStyle(h, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
		headerContainer.Add(label)
	}

	return headerContainer
}

func (pl *ProductList) showDeleteConfirmation(product models.Product) {
	dialog.ShowConfirm("Подтверждение удаления",
		fmt.Sprintf("Вы уверены, что хотите удалить товар '%s' (SKU: %s)?", product.Name, product.SKU),
		func(confirm bool) {
			if confirm {
				pl.deleteProduct(product)
			}
		},
		pl.mainWindow.window,
	)
}

func (pl *ProductList) deleteProduct(product models.Product) {
	err := pl.mainWindow.api.DeleteProduct(product.ID)
	if err != nil {
		dialog.ShowError(fmt.Errorf("Ошибка при удалении: %v", err), pl.mainWindow.window)
		return
	}

	pl.mainWindow.RefreshData()
	pl.mainWindow.statusBar.SetText(fmt.Sprintf("Товар '%s' удален", product.Name))
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
