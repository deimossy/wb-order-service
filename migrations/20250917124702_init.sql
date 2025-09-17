-- +goose Up
-- +goose StatementBegin

CREATE TABLE deliveries (
    delivery_uid UUID PRIMARY KEY,
    name TEXT NOT NULL,
    phone TEXT NOT NULL,
    zip TEXT NOT NULL,
    city TEXT NOT NULL,
    address TEXT NOT NULL,
    region TEXT NOT NULL,
    email TEXT NOT NULL
);

CREATE TABLE payments (
    payment_uid UUID PRIMARY KEY,
    transaction TEXT NOT NULL,
    request_id TEXT,
    currency TEXT NOT NULL,
    provider TEXT NOT NULL,
    amount INTEGER NOT NULL,
    payment_dt INTEGER NOT NULL,
    bank TEXT NOT NULL,
    delivery_cost INTEGER NOT NULL,
    goods_total INTEGER NOT NULL,
    custom_fee INTEGER NOT NULL
);

CREATE TABLE orders (
    order_uid UUID PRIMARY KEY,
    track_number TEXT NOT NULL,
    entry TEXT NOT NULL,
    delivery_uid UUID NOT NULL REFERENCES deliveries(delivery_uid) ON DELETE CASCADE,
    payment_uid UUID NOT NULL REFERENCES payments(payment_uid) ON DELETE CASCADE,
    locale TEXT NOT NULL,
    internal_signature TEXT,
    customer_id TEXT NOT NULL,
    delivery_service TEXT NOT NULL,
    shardkey TEXT NOT NULL,
    sm_id INTEGER NOT NULL,
    date_created TIMESTAMPTZ NOT NULL,
    oof_shard TEXT NOT NULL
);

CREATE TABLE items (
    item_uid UUID PRIMARY KEY,
    chrt_id INTEGER NOT NULL,
    track_number TEXT NOT NULL,
    rid TEXT NOT NULL,
    name TEXT NOT NULL,
    brand TEXT NOT NULL,
    size TEXT NOT NULL,
    nm_id INTEGER NOT NULL,
    status INTEGER NOT NULL
);

CREATE TABLE order_items (
    order_item_uid UUID PRIMARY KEY,
    order_uid UUID NOT NULL REFERENCES orders(order_uid) ON DELETE CASCADE,
    item_uid UUID NOT NULL REFERENCES items(item_uid) ON DELETE CASCADE,
    price INTEGER NOT NULL,
    sale INTEGER NOT NULL,
    total_price INTEGER NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 1
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS items;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS deliveries;

-- +goose StatementEnd
