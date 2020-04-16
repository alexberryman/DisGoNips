-- name: GetWeeksPriceHistoryByUser :many
SELECT *
FROM prices
WHERE discord_id = $1
  and day_of_year > extract(DOY FROM now()) - 7
  and year = extract(year from now())
order by day_of_year, meridiem;

-- name: ListPrices :many
SELECT *
FROM prices
ORDER BY created_at;

-- name: CountPricesByDiscordId :one
SELECT count(*)
FROM prices
where discord_id = $1;

-- name: CreatePrice :one
INSERT INTO prices (discord_id, price, meridiem, day_of_week, day_of_year, year)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: DeletePricesForUser :exec
DELETE
FROM prices
WHERE discord_id = $1;