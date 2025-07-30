-- name: CreateSubscription :exec
INSERT INTO subscriptions (id, service_name, price, user_id, start_date, end_date)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetSubscriptionByID :one
SELECT * FROM subscriptions WHERE id = $1;

-- name: ListSubscriptionsPaginated :many
SELECT *
FROM subscriptions
WHERE ($1::uuid IS NULL OR user_id = $1)
  AND ($2::text IS NULL OR service_name ILIKE '%' || $2 || '%')
ORDER BY start_date DESC
LIMIT $3 OFFSET $4;

-- name: UpdateSubscription :exec
UPDATE subscriptions
SET service_name = $2, price = $3, start_date = $4, end_date = $5
WHERE id = $1;

-- name: DeleteSubscription :exec
DELETE FROM subscriptions WHERE id = $1;

-- name: AggregateCost :one
SELECT COALESCE(SUM(price), 0) FROM subscriptions
WHERE ($1::uuid IS NULL OR user_id = $1)
  AND ($2::text IS NULL OR service_name = $2)
  AND start_date <= $3
  AND (end_date IS NULL OR end_date >= $4);
