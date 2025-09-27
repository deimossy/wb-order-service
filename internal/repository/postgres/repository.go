package postgres

import (
	"context"
	"database/sql"

	"github.com/deimossy/wb-order-service/internal/models"
	customerrors "github.com/deimossy/wb-order-service/pkg/custom_errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type OrderRepository interface {
	Save(ctx context.Context, o *models.Order) error
	GetByID(ctx context.Context, uid string) (*models.Order, error)
	GetLastN(ctx context.Context, n int) ([]*models.Order, error)
}

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Save(ctx context.Context, o *models.Order) error {
	if o.OrderUID == "" {
		o.OrderUID = uuid.New().String()
	}
	if o.Delivery != nil && o.Delivery.DeliveryUID == "" {
		o.Delivery.DeliveryUID = uuid.New().String()
	}
	if o.Payment != nil && o.Payment.PaymentUID == "" {
		o.Payment.PaymentUID = uuid.New().String()
	}
	for _, oi := range o.Items {
		if oi.OrderItemUID == "" {
			oi.OrderItemUID = uuid.New().String()
		}
		if oi.Item != nil && oi.Item.ItemUID == "" {
			oi.Item.ItemUID = uuid.New().String()
		}
		oi.OrderUID = o.OrderUID
		oi.ItemUID = oi.Item.ItemUID
	}

	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, deliveryInsertQuery,
		o.Delivery.DeliveryUID,
		o.Delivery.Name,
		o.Delivery.Phone,
		o.Delivery.Zip,
		o.Delivery.City,
		o.Delivery.Address,
		o.Delivery.Region,
		o.Delivery.Email); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, paymentInsertQuery,
		o.Payment.PaymentUID,
		o.Payment.Transaction,
		o.Payment.RequestID,
		o.Payment.Currency,
		o.Payment.Provider,
		o.Payment.Amount,
		o.Payment.PaymentDT,
		o.Payment.Bank,
		o.Payment.DeliveryCost,
		o.Payment.GoodsTotal,
		o.Payment.CustomFee); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, orderInsertQuery,
		o.OrderUID,
		o.TrackNumber,
		o.Entry,
		o.Delivery.DeliveryUID,
		o.Payment.PaymentUID,
		o.Locale,
		o.InternalSignature,
		o.CustomerID,
		o.DeliveryService,
		o.ShardKey,
		o.SmID,
		o.DateCreated,
		o.OofShard); err != nil {
		return err
	}

	for _, oi := range o.Items {
		if _, err := tx.ExecContext(ctx, itemInsertQuery,
			oi.Item.ItemUID,
			oi.Item.ChrtID,
			oi.Item.TrackNumber,
			oi.Item.RID,
			oi.Item.Name,
			oi.Item.Brand,
			oi.Item.Size,
			oi.Item.NmID,
			oi.Item.Status); err != nil {
			return err
		}

		if _, err := tx.ExecContext(ctx, orderItemInsertQuery,
			oi.OrderItemUID,
			o.OrderUID,
			oi.ItemUID,
			oi.Price,
			oi.Sale,
			oi.TotalPrice,
			oi.Quantity); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetByID(ctx context.Context, uid string) (*models.Order, error) {
	var o models.Order
	if err := r.db.GetContext(ctx, &o, getOrderByID, uid); err != nil {
		return nil, customerrors.ErrOrderNotFound
	}

	var d models.Delivery
	if err := r.db.GetContext(ctx, &d, getDeliveryByUID, o.DeliveryUID); err == nil {
		o.Delivery = &d
	}

	var p models.Payment
	if err := r.db.GetContext(ctx, &p, getPaymentByUID, o.PaymentUID); err == nil {
		o.Payment = &p
	}

	rows, err := r.db.QueryContext(ctx, getItemsByOrderUID, uid)
	if err != nil {
		return &o, nil
	}
	defer rows.Close()

	var items []*models.OrderItem
	for rows.Next() {
		var (
			it models.Item
			oi models.OrderItem
		)
		if err = rows.Scan(
			&it.ItemUID,
			&it.ChrtID,
			&it.TrackNumber,
			&it.RID,
			&it.Name,
			&it.Brand,
			&it.Size,
			&it.NmID,
			&it.Status,
			&oi.OrderItemUID,
			&oi.Price,
			&oi.Sale,
			&oi.TotalPrice,
			&oi.Quantity); err != nil {
			continue
		}

		oi.ItemUID = it.ItemUID
		oi.Item = &it
		oi.OrderUID = o.OrderUID
		items = append(items, &oi)
	}

	o.Items = items

	return &o, nil
}

func (r *Repository) GetLastN(ctx context.Context, n int) ([]*models.Order, error) {
	type row struct {
		OrderUID string `db:"order_uid"`
	}
	var rows []row
	if err := r.db.SelectContext(ctx, &rows, getLastNOrders, n); err != nil {
		return nil, err
	}

	var out []*models.Order
	for _, rrow := range rows {
		o, err := r.GetByID(ctx, rrow.OrderUID)
		if err == nil {
			out = append(out, o)
		}
	}

	return out, nil
}
