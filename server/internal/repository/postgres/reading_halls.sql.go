// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: reading_halls.sql

package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/govalues/decimal"
)

const createReadingHall = `-- name: CreateReadingHall :one
INSERT INTO reading_halls (hall_name, specialization, total_seats)
VALUES ($1, $2, $3)
RETURNING id, hall_name, specialization, total_seats, current_visitors
`

type CreateReadingHallParams struct {
	HallName       string  `json:"hall_name"`
	Specialization *string `json:"specialization"`
	TotalSeats     int     `json:"total_seats"`
}

type CreateReadingHallRow struct {
	ID              uuid.UUID `json:"id"`
	HallName        string    `json:"hall_name"`
	Specialization  *string   `json:"specialization"`
	TotalSeats      int       `json:"total_seats"`
	CurrentVisitors *int      `json:"current_visitors"`
}

func (q *Queries) CreateReadingHall(ctx context.Context, arg CreateReadingHallParams) (*CreateReadingHallRow, error) {
	row := q.db.QueryRow(ctx, createReadingHall, arg.HallName, arg.Specialization, arg.TotalSeats)
	var i CreateReadingHallRow
	err := row.Scan(
		&i.ID,
		&i.HallName,
		&i.Specialization,
		&i.TotalSeats,
		&i.CurrentVisitors,
	)
	return &i, err
}

const getAllReadingHalls = `-- name: GetAllReadingHalls :many
SELECT id, hall_name, specialization, total_seats, current_visitors
FROM reading_halls
ORDER BY hall_name
`

type GetAllReadingHallsRow struct {
	ID              uuid.UUID `json:"id"`
	HallName        string    `json:"hall_name"`
	Specialization  *string   `json:"specialization"`
	TotalSeats      int       `json:"total_seats"`
	CurrentVisitors *int      `json:"current_visitors"`
}

func (q *Queries) GetAllReadingHalls(ctx context.Context) ([]*GetAllReadingHallsRow, error) {
	rows, err := q.db.Query(ctx, getAllReadingHalls)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetAllReadingHallsRow{}
	for rows.Next() {
		var i GetAllReadingHallsRow
		if err := rows.Scan(
			&i.ID,
			&i.HallName,
			&i.Specialization,
			&i.TotalSeats,
			&i.CurrentVisitors,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getHallsDashboard = `-- name: GetHallsDashboard :many
SELECT
    rh.id,
    rh.hall_name,
    rh.specialization,
    rh.total_seats,
    rh.current_visitors,
    ROUND(
        (rh.current_visitors::numeric / rh.total_seats * 100), 2
    ) as occupancy_percentage,
    (rh.total_seats - rh.current_visitors) as free_seats,
    COALESCE(daily_stats.visits_today, 0) as visits_today,
    COALESCE(daily_stats.unique_visitors_today, 0) as unique_visitors_today
FROM reading_halls rh
LEFT JOIN (
    SELECT
        hall_id,
        COUNT(*) as visits_today,
        COUNT(DISTINCT reader_id) as unique_visitors_today
    FROM hall_visits
    WHERE DATE(visit_time) = CURRENT_DATE
    GROUP BY hall_id
) daily_stats ON rh.id = daily_stats.hall_id
ORDER BY rh.hall_name
`

type GetHallsDashboardRow struct {
	ID                  uuid.UUID       `json:"id"`
	HallName            string          `json:"hall_name"`
	Specialization      *string         `json:"specialization"`
	TotalSeats          int             `json:"total_seats"`
	CurrentVisitors     *int            `json:"current_visitors"`
	OccupancyPercentage decimal.Decimal `json:"occupancy_percentage"`
	FreeSeats           int32           `json:"free_seats"`
	VisitsToday         int64           `json:"visits_today"`
	UniqueVisitorsToday int64           `json:"unique_visitors_today"`
}

func (q *Queries) GetHallsDashboard(ctx context.Context) ([]*GetHallsDashboardRow, error) {
	rows, err := q.db.Query(ctx, getHallsDashboard)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetHallsDashboardRow{}
	for rows.Next() {
		var i GetHallsDashboardRow
		if err := rows.Scan(
			&i.ID,
			&i.HallName,
			&i.Specialization,
			&i.TotalSeats,
			&i.CurrentVisitors,
			&i.OccupancyPercentage,
			&i.FreeSeats,
			&i.VisitsToday,
			&i.UniqueVisitorsToday,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getReadingHallById = `-- name: GetReadingHallById :one
SELECT id, hall_name, specialization, total_seats, current_visitors
FROM reading_halls
WHERE id = $1
`

type GetReadingHallByIdRow struct {
	ID              uuid.UUID `json:"id"`
	HallName        string    `json:"hall_name"`
	Specialization  *string   `json:"specialization"`
	TotalSeats      int       `json:"total_seats"`
	CurrentVisitors *int      `json:"current_visitors"`
}

func (q *Queries) GetReadingHallById(ctx context.Context, id uuid.UUID) (*GetReadingHallByIdRow, error) {
	row := q.db.QueryRow(ctx, getReadingHallById, id)
	var i GetReadingHallByIdRow
	err := row.Scan(
		&i.ID,
		&i.HallName,
		&i.Specialization,
		&i.TotalSeats,
		&i.CurrentVisitors,
	)
	return &i, err
}

const updateHallVisitorCount = `-- name: UpdateHallVisitorCount :exec
UPDATE reading_halls
SET current_visitors = current_visitors + $1
WHERE id = $2
`

type UpdateHallVisitorCountParams struct {
	Change *int      `json:"change"`
	HallID uuid.UUID `json:"hall_id"`
}

func (q *Queries) UpdateHallVisitorCount(ctx context.Context, arg UpdateHallVisitorCountParams) error {
	_, err := q.db.Exec(ctx, updateHallVisitorCount, arg.Change, arg.HallID)
	return err
}

const updateReadingHall = `-- name: UpdateReadingHall :one
UPDATE reading_halls
SET hall_name = $1, specialization = $2, total_seats = $3
WHERE id = $4
RETURNING id, hall_name, specialization, total_seats, current_visitors
`

type UpdateReadingHallParams struct {
	HallName       string    `json:"hall_name"`
	Specialization *string   `json:"specialization"`
	TotalSeats     int       `json:"total_seats"`
	ID             uuid.UUID `json:"id"`
}

type UpdateReadingHallRow struct {
	ID              uuid.UUID `json:"id"`
	HallName        string    `json:"hall_name"`
	Specialization  *string   `json:"specialization"`
	TotalSeats      int       `json:"total_seats"`
	CurrentVisitors *int      `json:"current_visitors"`
}

func (q *Queries) UpdateReadingHall(ctx context.Context, arg UpdateReadingHallParams) (*UpdateReadingHallRow, error) {
	row := q.db.QueryRow(ctx, updateReadingHall,
		arg.HallName,
		arg.Specialization,
		arg.TotalSeats,
		arg.ID,
	)
	var i UpdateReadingHallRow
	err := row.Scan(
		&i.ID,
		&i.HallName,
		&i.Specialization,
		&i.TotalSeats,
		&i.CurrentVisitors,
	)
	return &i, err
}
