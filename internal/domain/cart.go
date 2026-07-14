package domain

type Cart struct {
	Subtotal      float64
	DiscountTotal float64
	Total         float64
	Currency      string
	Items         []CartItem
}

type CartItem struct {
	ProductID string
	Quantity  int
	UnitPrice float64
}
