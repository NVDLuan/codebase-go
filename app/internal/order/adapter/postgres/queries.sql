-- name: UpsertOrder :exec
INSERT INTO orders (id, customer_id, status)
VALUES ($1, $2, $3)
ON CONFLICT (id) DO UPDATE SET status = EXCLUDED.status, updated_at = now();

-- name: DeleteOrderItems :exec
DELETE FROM order_items WHERE order_id = $1;

-- name: InsertOrderItem :exec
INSERT INTO order_items (order_id, position, product_id, unit_price, currency, quantity)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetOrder :one
SELECT id, customer_id, status FROM orders WHERE id = $1;

-- name: ListOrderItems :many
SELECT product_id, unit_price, currency, quantity
FROM order_items
WHERE order_id = $1
ORDER BY position;
