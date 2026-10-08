-- +goose Up
CREATE TABLE orders (
    id          UUID PRIMARY KEY,
    customer_id UUID        NOT NULL,
    status      TEXT        NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE order_items (
    order_id   UUID    NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
    position   INTEGER NOT NULL,
    product_id UUID    NOT NULL,
    unit_price BIGINT  NOT NULL,
    currency   TEXT    NOT NULL,
    quantity   INTEGER NOT NULL CHECK (quantity > 0),
    PRIMARY KEY (order_id, position)
);

CREATE INDEX idx_orders_customer_id ON orders (customer_id);

-- +goose Down
DROP TABLE order_items;
DROP TABLE orders;
