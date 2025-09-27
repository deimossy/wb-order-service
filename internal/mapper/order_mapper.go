package mapper

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/deimossy/wb-order-service/internal/dto"
	"github.com/deimossy/wb-order-service/internal/models"
	"github.com/google/uuid"
)

func MapJSONToModel(data []byte) (*models.Order, error) {
	var dtoOrder dto.OrderDTO
	if err := json.Unmarshal(data, &dtoOrder); err != nil {
		return nil, err
	}

	return MapDTOToModel(&dtoOrder)
}

func MapDTOToModel(d *dto.OrderDTO) (*models.Order, error) {
	if d == nil {
		return nil, errors.New("nil dto")
	}

	ou := d.OrderUID
	du := uuid.New().String()
	pu := uuid.New().String()

	dateCreated, err := time.Parse(time.RFC3339, d.DateCreated)
	if err != nil {
		return nil, err
	}

	order := &models.Order{
		OrderUID:          ou,
		TrackNumber:       d.TrackNumber,
		Entry:             d.Entry,
		Locale:            d.Locale,
		InternalSignature: d.InternalSignature,
		CustomerID:        d.CustomerID,
		DeliveryService:   d.DeliveryService,
		ShardKey:          d.ShardKey,
		SmID:              d.SmID,
		DateCreated:       dateCreated,
		OofShard:          d.OofShard,
		DeliveryUID:       du,
		PaymentUID:        pu,
	}

	order.Delivery = &models.Delivery{
		DeliveryUID: du,
		Name:        d.Delivery.Name,
		Phone:       d.Delivery.Phone,
		Zip:         d.Delivery.Zip,
		City:        d.Delivery.City,
		Address:     d.Delivery.Address,
		Region:      d.Delivery.Region,
		Email:       d.Delivery.Email,
	}

	order.Payment = &models.Payment{
		PaymentUID:   pu,
		Transaction:  d.Payment.Transaction,
		RequestID:    d.Payment.RequestID,
		Currency:     d.Payment.Currency,
		Provider:     d.Payment.Provider,
		Amount:       d.Payment.Amount,
		PaymentDT:    d.Payment.PaymentDT,
		Bank:         d.Payment.Bank,
		DeliveryCost: d.Payment.DeliveryCost,
		GoodsTotal:   d.Payment.GoodsTotal,
		CustomFee:    d.Payment.CustomFee,
	}

	for _, it := range d.Items {
		itemUID := uuid.New().String()
		orderItemUID := uuid.New().String()

		item := &models.Item{
			ItemUID:     itemUID,
			ChrtID:      it.ChrtID,
			TrackNumber: it.TrackNumber,
			RID:         it.RID,
			Name:        it.Name,
			Brand:       it.Brand,
			Size:        it.Size,
			NmID:        it.NmID,
			Status:      it.Status,
		}

		orderItem := &models.OrderItem{
			OrderItemUID: orderItemUID,
			OrderUID:     order.OrderUID,
			ItemUID:      itemUID,
			Item:         item,
			Price:        it.Price,
			Sale:         it.Sale,
			TotalPrice:   it.TotalPrice,
			Quantity:     1,
		}

		order.Items = append(order.Items, orderItem)
	}

	return order, nil
}
