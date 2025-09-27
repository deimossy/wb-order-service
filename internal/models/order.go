package models

import (
	"time"
)

type Order struct {
	OrderUID          string    `db:"order_uid"`
	TrackNumber       string    `db:"track_number"`
	Entry             string    `db:"entry"`
	DeliveryUID       string    `db:"delivery_uid"`
	PaymentUID        string    `db:"payment_uid"`
	Locale            string    `db:"locale"`
	InternalSignature string    `db:"internal_signature"`
	CustomerID        string    `db:"customer_id"`
	DeliveryService   string    `db:"delivery_service"`
	ShardKey          string    `db:"shardkey"`
	SmID              int       `db:"sm_id"`
	DateCreated       time.Time `db:"date_created"`
	OofShard          string    `db:"oof_shard"`

	Delivery *Delivery    `db:"-"`
	Payment  *Payment     `db:"-"`
	Items    []*OrderItem `db:"-"`
}
