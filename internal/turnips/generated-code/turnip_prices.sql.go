// Code generated by sqlc. DO NOT EDIT.
// source: turnip_prices.sql

package turnips

import (
	"context"
	"database/sql"
	"time"
)

const countPricesByDiscordId = `-- name: CountPricesByDiscordId :one
SELECT count(*)
FROM turnip_prices
where discord_id = $1
`

func (q *Queries) CountPricesByDiscordId(ctx context.Context, discordID string) (int64, error) {
	row := q.queryRow(ctx, q.countPricesByDiscordIdStmt, countPricesByDiscordId, discordID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createPrice = `-- name: CreatePrice :one
INSERT INTO turnip_prices (discord_id, price, am_pm, day_of_week, day_of_year, year)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, discord_id, price, am_pm, day_of_week, day_of_year, year, created_at, week
`

type CreatePriceParams struct {
	DiscordID string `json:"discord_id"`
	Price     int32  `json:"price"`
	AmPm      AmPm   `json:"am_pm"`
	DayOfWeek int32  `json:"day_of_week"`
	DayOfYear int32  `json:"day_of_year"`
	Year      int32  `json:"year"`
}

func (q *Queries) CreatePrice(ctx context.Context, arg CreatePriceParams) (TurnipPrice, error) {
	row := q.queryRow(ctx, q.createPriceStmt, createPrice,
		arg.DiscordID,
		arg.Price,
		arg.AmPm,
		arg.DayOfWeek,
		arg.DayOfYear,
		arg.Year,
	)
	var i TurnipPrice
	err := row.Scan(
		&i.ID,
		&i.DiscordID,
		&i.Price,
		&i.AmPm,
		&i.DayOfWeek,
		&i.DayOfYear,
		&i.Year,
		&i.CreatedAt,
		&i.Week,
	)
	return i, err
}

const deletePricesForUser = `-- name: DeletePricesForUser :exec
DELETE
FROM turnip_prices
WHERE discord_id = $1
`

func (q *Queries) DeletePricesForUser(ctx context.Context, discordID string) error {
	_, err := q.exec(ctx, q.deletePricesForUserStmt, deletePricesForUser, discordID)
	return err
}

const getLastWeeksPriceHistoryByServer = `-- name: GetLastWeeksPriceHistoryByServer :many
SELECT id, turnip_prices.discord_id, price, am_pm, day_of_week, day_of_year, year, created_at, week, nick.discord_id, server_id, nickname
FROM turnip_prices
         inner join nicknames nick on turnip_prices.discord_id = nick.discord_id
WHERE nick.server_id = $1
  and year = extract(year from now())
  and week = extract(week from now()) - 1
order by day_of_year, am_pm
`

type GetLastWeeksPriceHistoryByServerRow struct {
	ID          int64         `json:"id"`
	DiscordID   string        `json:"discord_id"`
	Price       int32         `json:"price"`
	AmPm        AmPm          `json:"am_pm"`
	DayOfWeek   int32         `json:"day_of_week"`
	DayOfYear   int32         `json:"day_of_year"`
	Year        int32         `json:"year"`
	CreatedAt   time.Time     `json:"created_at"`
	Week        sql.NullInt32 `json:"week"`
	DiscordID_2 string        `json:"discord_id_2"`
	ServerID    string        `json:"server_id"`
	Nickname    string        `json:"nickname"`
}

func (q *Queries) GetLastWeeksPriceHistoryByServer(ctx context.Context, serverID string) ([]GetLastWeeksPriceHistoryByServerRow, error) {
	rows, err := q.query(ctx, q.getLastWeeksPriceHistoryByServerStmt, getLastWeeksPriceHistoryByServer, serverID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetLastWeeksPriceHistoryByServerRow
	for rows.Next() {
		var i GetLastWeeksPriceHistoryByServerRow
		if err := rows.Scan(
			&i.ID,
			&i.DiscordID,
			&i.Price,
			&i.AmPm,
			&i.DayOfWeek,
			&i.DayOfYear,
			&i.Year,
			&i.CreatedAt,
			&i.Week,
			&i.DiscordID_2,
			&i.ServerID,
			&i.Nickname,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getWeeksPriceHistoryByAccount = `-- name: GetWeeksPriceHistoryByAccount :many
SELECT id, discord_id, price, am_pm, day_of_week, day_of_year, year, created_at, week
FROM turnip_prices
WHERE discord_id = $1
  and day_of_year > extract(DOY FROM now()) - 7
  and year = extract(year from now())
order by day_of_year, am_pm
`

func (q *Queries) GetWeeksPriceHistoryByAccount(ctx context.Context, discordID string) ([]TurnipPrice, error) {
	rows, err := q.query(ctx, q.getWeeksPriceHistoryByAccountStmt, getWeeksPriceHistoryByAccount, discordID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []TurnipPrice
	for rows.Next() {
		var i TurnipPrice
		if err := rows.Scan(
			&i.ID,
			&i.DiscordID,
			&i.Price,
			&i.AmPm,
			&i.DayOfWeek,
			&i.DayOfYear,
			&i.Year,
			&i.CreatedAt,
			&i.Week,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getWeeksPriceHistoryByServer = `-- name: GetWeeksPriceHistoryByServer :many
SELECT id, turnip_prices.discord_id, price, am_pm, day_of_week, day_of_year, year, created_at, week, nick.discord_id, server_id, nickname
FROM turnip_prices
         inner join nicknames nick on turnip_prices.discord_id = nick.discord_id
WHERE nick.server_id = $1
  and year = extract(year from now())
  and week = extract(week from now())
order by day_of_year, am_pm
`

type GetWeeksPriceHistoryByServerRow struct {
	ID          int64         `json:"id"`
	DiscordID   string        `json:"discord_id"`
	Price       int32         `json:"price"`
	AmPm        AmPm          `json:"am_pm"`
	DayOfWeek   int32         `json:"day_of_week"`
	DayOfYear   int32         `json:"day_of_year"`
	Year        int32         `json:"year"`
	CreatedAt   time.Time     `json:"created_at"`
	Week        sql.NullInt32 `json:"week"`
	DiscordID_2 string        `json:"discord_id_2"`
	ServerID    string        `json:"server_id"`
	Nickname    string        `json:"nickname"`
}

func (q *Queries) GetWeeksPriceHistoryByServer(ctx context.Context, serverID string) ([]GetWeeksPriceHistoryByServerRow, error) {
	rows, err := q.query(ctx, q.getWeeksPriceHistoryByServerStmt, getWeeksPriceHistoryByServer, serverID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetWeeksPriceHistoryByServerRow
	for rows.Next() {
		var i GetWeeksPriceHistoryByServerRow
		if err := rows.Scan(
			&i.ID,
			&i.DiscordID,
			&i.Price,
			&i.AmPm,
			&i.DayOfWeek,
			&i.DayOfYear,
			&i.Year,
			&i.CreatedAt,
			&i.Week,
			&i.DiscordID_2,
			&i.ServerID,
			&i.Nickname,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listPrices = `-- name: ListPrices :many
SELECT id, discord_id, price, am_pm, day_of_week, day_of_year, year, created_at, week
FROM turnip_prices
ORDER BY created_at
`

func (q *Queries) ListPrices(ctx context.Context) ([]TurnipPrice, error) {
	rows, err := q.query(ctx, q.listPricesStmt, listPrices)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []TurnipPrice
	for rows.Next() {
		var i TurnipPrice
		if err := rows.Scan(
			&i.ID,
			&i.DiscordID,
			&i.Price,
			&i.AmPm,
			&i.DayOfWeek,
			&i.DayOfYear,
			&i.Year,
			&i.CreatedAt,
			&i.Week,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updatePrice = `-- name: UpdatePrice :one
update turnip_prices
set price = $2
where discord_id = $1
  and am_pm = $3
  and day_of_week = $4
  and day_of_year = $5
  and year = $6
returning id, discord_id, price, am_pm, day_of_week, day_of_year, year, created_at, week
`

type UpdatePriceParams struct {
	DiscordID string `json:"discord_id"`
	Price     int32  `json:"price"`
	AmPm      AmPm   `json:"am_pm"`
	DayOfWeek int32  `json:"day_of_week"`
	DayOfYear int32  `json:"day_of_year"`
	Year      int32  `json:"year"`
}

func (q *Queries) UpdatePrice(ctx context.Context, arg UpdatePriceParams) (TurnipPrice, error) {
	row := q.queryRow(ctx, q.updatePriceStmt, updatePrice,
		arg.DiscordID,
		arg.Price,
		arg.AmPm,
		arg.DayOfWeek,
		arg.DayOfYear,
		arg.Year,
	)
	var i TurnipPrice
	err := row.Scan(
		&i.ID,
		&i.DiscordID,
		&i.Price,
		&i.AmPm,
		&i.DayOfWeek,
		&i.DayOfYear,
		&i.Year,
		&i.CreatedAt,
		&i.Week,
	)
	return i, err
}
