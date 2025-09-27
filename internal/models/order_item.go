package models

type OrderItem struct {
	OrderItemUID string `db:"-"`
	OrderUID     string `db:"-"`
	ItemUID      string `db:"-"`

	Item       *Item `db:"-"`
	Price      int   `db:"price"`
	Sale       int   `db:"sale"`
	TotalPrice int   `db:"total_price"`
	Quantity   int   `db:"quantity"`
}
