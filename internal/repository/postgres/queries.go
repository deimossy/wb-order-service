package postgres

const (
	deliveryInsertQuery = `
	INSERT INTO deliveries (delivery_uid, name, phone, zip, city, address, region, email)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	ON CONFLICT (delivery_uid) DO NOTHING
	`

	paymentInsertQuery = `
	INSERT INTO payments (payment_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
	ON CONFLICT (payment_uid) DO NOTHING
	`

	orderInsertQuery = `
	INSERT INTO orders (order_uid, track_number, entry, delivery_uid, payment_uid, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
	ON CONFLICT (order_uid) DO NOTHING
	`

	itemInsertQuery = `
	INSERT INTO items (item_uid, chrt_id, track_number, rid, name, brand, size, nm_id, status)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	ON CONFLICT (item_uid) DO NOTHING
	`

	orderItemInsertQuery = `
	INSERT INTO order_items (order_item_uid, order_uid, item_uid, price, sale, total_price, quantity)
	VALUES ($1,$2,$3,$4,$5,$6,$7)
	ON CONFLICT (order_item_uid) DO NOTHING
	`

	getOrderByID = `
	SELECT
		o.order_uid, o.track_number, o.entry, o.delivery_uid, o.payment_uid, o.locale,
		o.internal_signature, o.customer_id, o.delivery_service, o.shardkey, o.sm_id, o.date_created, o.oof_shard
	FROM orders o
	WHERE o.order_uid = $1
	`

	getDeliveryByUID = `
	SELECT delivery_uid, name, phone, zip, city, address, region, email
	FROM deliveries
	WHERE delivery_uid = $1
	`

	getPaymentByUID = `
	SELECT payment_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
	FROM payments
	WHERE payment_uid = $1
	`

	getItemsByOrderUID = `
	SELECT
		i.item_uid, i.chrt_id, i.track_number, i.rid, i.name, i.brand, i.size, i.nm_id, i.status,
		oi.order_item_uid, oi.price, oi.sale, oi.total_price, oi.quantity
	FROM order_items oi
	JOIN items i ON oi.item_uid = i.item_uid
	WHERE oi.order_uid = $1
	ORDER BY oi.order_item_uid
	`

	getLastNOrders = `
	SELECT order_uid 
	FROM orders 
	ORDER BY date_created DESC LIMIT $1
	`
)