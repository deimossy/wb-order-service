package models

type Delivery struct {
	DeliveryUID string `db:"-"`
	Name        string `db:"name"`
	Phone       string `db:"phone"`
	Zip         string `db:"zip"`
	City        string `db:"city"`
	Address     string `db:"address"`
	Region      string `db:"region"`
	Email       string `db:"email"`
}
