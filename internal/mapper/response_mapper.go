package mapper

import (
	"time"

	"github.com/deimossy/wb-order-service/internal/models"
)

type DeliveryResp struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type PaymentResp struct {
	Transaction  string `json:"transaction"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDT    int64  `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

type ItemResp struct {
	ChrtID      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	RID         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmID        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

type OrderResp struct {
	OrderUID          string       `json:"order_uid"`
	TrackNumber       string       `json:"track_number"`
	Entry             string       `json:"entry"`
	Delivery          DeliveryResp `json:"delivery"`
	Payment           PaymentResp  `json:"payment"`
	Items             []ItemResp   `json:"items"`
	Locale            string       `json:"locale"`
	InternalSignature string       `json:"internal_signature"`
	CustomerID        string       `json:"customer_id"`
	DeliveryService   string       `json:"delivery_service"`
	ShardKey          string       `json:"shardkey"`
	SmID              int          `json:"sm_id"`
	DateCreated       string       `json:"date_created"`
	OofShard          string       `json:"oof_shard"`
}

func ModelToResponse(o *models.Order) *OrderResp {
	if o == nil {
		return nil
	}

	dr := DeliveryResp{}
	if o.Delivery != nil {
		dr = DeliveryResp{
			Name:    o.Delivery.Name,
			Phone:   o.Delivery.Phone,
			Zip:     o.Delivery.Zip,
			City:    o.Delivery.City,
			Address: o.Delivery.Address,
			Region:  o.Delivery.Region,
			Email:   o.Delivery.Email,
		}
	}

	pr := PaymentResp{}
	if o.Payment != nil {
		pr = PaymentResp{
			Transaction:  o.Payment.Transaction,
			RequestID:    o.Payment.RequestID,
			Currency:     o.Payment.Currency,
			Provider:     o.Payment.Provider,
			Amount:       o.Payment.Amount,
			PaymentDT:    o.Payment.PaymentDT,
			Bank:         o.Payment.Bank,
			DeliveryCost: o.Payment.DeliveryCost,
			GoodsTotal:   o.Payment.GoodsTotal,
			CustomFee:    o.Payment.CustomFee,
		}
	}

	var items []ItemResp
	for _, oi := range o.Items {
		if oi == nil || oi.Item == nil {
			continue
		}
		items = append(items, ItemResp{
			ChrtID:      oi.Item.ChrtID,
			TrackNumber: oi.Item.TrackNumber,
			Price:       oi.Price,
			RID:         oi.Item.RID,
			Name:        oi.Item.Name,
			Sale:        oi.Sale,
			Size:        oi.Item.Size,
			TotalPrice:  oi.TotalPrice,
			NmID:        oi.Item.NmID,
			Brand:       oi.Item.Brand,
			Status:      oi.Item.Status,
		})
	}

	return &OrderResp{
		OrderUID:          o.OrderUID,
		TrackNumber:       o.TrackNumber,
		Entry:             o.Entry,
		Delivery:          dr,
		Payment:           pr,
		Items:             items,
		Locale:            o.Locale,
		InternalSignature: o.InternalSignature,
		CustomerID:        o.CustomerID,
		DeliveryService:   o.DeliveryService,
		ShardKey:          o.ShardKey,
		SmID:              o.SmID,
		DateCreated:       o.DateCreated.Format(time.RFC3339),
		OofShard:          o.OofShard,
	}
}
