package gui

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"SanWarehouse/models"
)

type ProductForm struct {
	window  fyne.Window
	product *models.Product
	onSave  func(*models.Product)

	// Поля формы
	skuLabel         *widget.Label
	skuEntry         *widget.Entry
	nameEntry        *widget.Entry
	categoryEntry    *widget.Entry
	brandEntry       *widget.Entry
	descEntry        *widget.Entry
	quantityEntry    *widget.Entry
	reservedEntry    *widget.Entry
	purchaseEntry    *widget.Entry
	sellingEntry     *widget.Entry
	minStockEntry    *widget.Entry
	locationEntry    *widget.Entry
	weightEntry      *widget.Entry
	dimensionsEntry  *widget.Entry
	materialEntry    *widget.Entry
	marketplaceEntry *widget.Entry
	activeCheck      *widget.Check
}

func NewProductForm(parent fyne.Window, product *models.Product, onSave func(*models.Product)) *ProductForm {
	pf := &ProductForm{
		window:  parent,
		product: product,
		onSave:  onSave,
	}

	pf.initFields()
	return pf
}

func (pf *ProductForm) initFields() {
	// Создаем поля ввода
	pf.skuEntry = widget.NewEntry()
	pf.nameEntry = widget.NewEntry()
	pf.categoryEntry = widget.NewEntry()
	pf.brandEntry = widget.NewEntry()
	pf.descEntry = widget.NewEntry()
	pf.quantityEntry = widget.NewEntry()
	pf.reservedEntry = widget.NewEntry()
	pf.purchaseEntry = widget.NewEntry()
	pf.sellingEntry = widget.NewEntry()
	pf.minStockEntry = widget.NewEntry()
	pf.locationEntry = widget.NewEntry()
	pf.weightEntry = widget.NewEntry()
	pf.dimensionsEntry = widget.NewEntry()
	pf.materialEntry = widget.NewEntry()
	pf.marketplaceEntry = widget.NewEntry()
	pf.activeCheck = widget.NewCheck("Активен", nil)

	// Если редактируем существующий товар, заполняем поля
	if pf.product != nil {
		pf.skuEntry.SetText(pf.product.SKU)
		pf.nameEntry.SetText(pf.product.Name)
		pf.categoryEntry.SetText(pf.product.Category)
		pf.brandEntry.SetText(pf.product.Brand)
		pf.descEntry.SetText(pf.product.Description)
		pf.quantityEntry.SetText(strconv.Itoa(pf.product.Quantity))
		pf.reservedEntry.SetText(strconv.Itoa(pf.product.ReservedQuantity))
		pf.purchaseEntry.SetText(strconv.FormatFloat(pf.product.PurchasePrice, 'f', 2, 64))
		pf.sellingEntry.SetText(strconv.FormatFloat(pf.product.SellingPrice, 'f', 2, 64))
		pf.minStockEntry.SetText(strconv.Itoa(pf.product.MinStockLevel))
		pf.locationEntry.SetText(pf.product.Location)
		pf.weightEntry.SetText(strconv.FormatFloat(pf.product.Weight, 'f', 2, 64))
		pf.dimensionsEntry.SetText(pf.product.Dimensions)
		pf.materialEntry.SetText(pf.product.Material)
		pf.marketplaceEntry.SetText(pf.product.MarketplaceID)
		pf.activeCheck.SetChecked(pf.product.IsActive)
	} else {
		pf.activeCheck.SetChecked(true)
	}
}

func (pf *ProductForm) Show() {
	// Создаем форму
	items := []*widget.FormItem{
		widget.NewFormItem("SKU*", pf.skuEntry),
		widget.NewFormItem("Наименование*", pf.nameEntry),
		widget.NewFormItem("Категория", pf.categoryEntry),
		widget.NewFormItem("Бренд", pf.brandEntry),
		widget.NewFormItem("Описание", pf.descEntry),
		widget.NewFormItem("Количество", pf.quantityEntry),
		widget.NewFormItem("Зарезервировано", pf.reservedEntry),
		widget.NewFormItem("Закупочная цена", pf.purchaseEntry),
		widget.NewFormItem("Цена продажи", pf.sellingEntry),
		widget.NewFormItem("Мин. уровень", pf.minStockEntry),
		widget.NewFormItem("Расположение", pf.locationEntry),
		widget.NewFormItem("Вес (кг)", pf.weightEntry),
		widget.NewFormItem("Габариты", pf.dimensionsEntry),
		widget.NewFormItem("Материал", pf.materialEntry),
		widget.NewFormItem("ID на маркетплейсе", pf.marketplaceEntry),
	}

	// Создаем контент с прокруткой
	content := container.NewVBox()

	for _, item := range items {
		content.Add(widget.NewLabel(item.Text))
		content.Add(container.NewPadded(
			item.Widget,
		))
	}

	content.Add(container.NewPadded(pf.activeCheck))

	scroll := container.NewScroll(content)
	scroll.SetMinSize(fyne.NewSize(500, 500))

	title := "Добавление товара"
	if pf.product != nil {
		title = "Редактирование товара"
	}

	dialog.ShowCustomConfirm(title, "Сохранить", "Отмена", scroll,
		func(confirm bool) {
			if confirm {
				pf.saveProduct()
			}
		},
		pf.window,
	)
}

func (pf *ProductForm) saveProduct() {
	// Валидация
	if pf.skuEntry.Text == "" || pf.nameEntry.Text == "" {
		err := pf.validationError("Заполните обязательные поля (SKU и Наименование)")
		dialog.ShowError(err, pf.window)
		return
	}

	// Создаем или обновляем продукт
	product := &models.Product{}
	if pf.product != nil {
		product = pf.product
	}

	product.SKU = pf.skuEntry.Text
	product.Name = pf.nameEntry.Text
	product.Category = pf.categoryEntry.Text
	product.Brand = pf.brandEntry.Text
	product.Description = pf.descEntry.Text

	// Конвертируем строки в числа
	quantity, _ := strconv.Atoi(pf.quantityEntry.Text)
	product.Quantity = quantity

	reserved, _ := strconv.Atoi(pf.reservedEntry.Text)
	product.ReservedQuantity = reserved

	purchase, _ := strconv.ParseFloat(pf.purchaseEntry.Text, 64)
	product.PurchasePrice = purchase

	selling, _ := strconv.ParseFloat(pf.sellingEntry.Text, 64)
	product.SellingPrice = selling

	minStock, _ := strconv.Atoi(pf.minStockEntry.Text)
	if minStock > 0 {
		product.MinStockLevel = minStock
	}

	product.Location = pf.locationEntry.Text

	weight, _ := strconv.ParseFloat(pf.weightEntry.Text, 64)
	product.Weight = weight

	product.Dimensions = pf.dimensionsEntry.Text
	product.Material = pf.materialEntry.Text
	product.MarketplaceID = pf.marketplaceEntry.Text
	product.IsActive = pf.activeCheck.Checked

	// Вызываем колбэк сохранения
	pf.onSave(product)
}

// ИСПРАВЛЕНО: Добавлен метод для создания error
func (pf *ProductForm) validationError(message string) error {
	return &validationError{message}
}

// ИСПРАВЛЕНО: Структура для ошибки валидации
type validationError struct {
	msg string
}

func (e *validationError) Error() string {
	return e.msg
}
