package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"github.com/maYkiss56/subscription-aggregation-service/internal/domain"
	"github.com/maYkiss56/subscription-aggregation-service/pkg/client/postgresql"
)

type SubRepository struct {
	pg *postgresql.PostgresClient
}

func New(pg *postgresql.PostgresClient) *SubRepository {
	return &SubRepository{pg: pg}
}

func (r *SubRepository) CreateSub(ctx context.Context, sub *domain.Sub) (id uuid.UUID, err error) {
	conn, err := r.pg.GetConnection(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to get connection: %w", err)
	}
	defer conn.Release()

	query := `
		insert into subscriptions
		(id, service_name, price, user_id, start_date, end_date)
		values ($1, $2, $3, $4, $5, $6)
		returning id
	`

	err = conn.QueryRow(
		ctx, query,
		sub.ID,
		sub.ServiceName,
		sub.Price,
		sub.UserID,
		sub.StartDate,
		sub.EndDate,
	).Scan(&sub.ID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create subsciption: %w", err)
	}

	return sub.ID, err
}

func (r *SubRepository) GetAllSubs(ctx context.Context) ([]*domain.Sub, error) {
	conn, err := r.pg.GetConnection(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed go get connection: %w", err)
	}
	defer conn.Release()

	query := `select
		id, service_name,
		price, user_id,
		start_date, end_date
		from subscriptions
	`

	rows, err := conn.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query subs: %w", err)
	}
	defer rows.Close()

	var subs []*domain.Sub
	for rows.Next() {
		var sub domain.Sub
		err := rows.Scan(
			&sub.ID,
			&sub.ServiceName,
			&sub.Price,
			&sub.UserID,
			&sub.StartDate,
			&sub.EndDate,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to san row subs: %w", err)
		}
		subs = append(subs, &sub)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return subs, nil
}

func (r *SubRepository) GetSubByUserID(ctx context.Context, userUID uuid.UUID) ([]*domain.Sub, error) {
	conn, err := r.pg.GetConnection(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection: %w", err)
	}
	defer conn.Release()

	query := `select
		id, service_name,
		price, user_id,
		start_date, end_date
		from subscriptions
		where user_id=$1
	`

	rows, err := conn.Query(ctx, query, userUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get subs by id: %w", err)
	}
	defer rows.Close()

	var subs []*domain.Sub
	for rows.Next() {
		var sub domain.Sub
		err := rows.Scan(
			&sub.ID,
			&sub.ServiceName,
			&sub.Price,
			&sub.UserID,
			&sub.StartDate,
			&sub.EndDate,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan subs: %w", err)
		}
		subs = append(subs, &sub)
	}

	return subs, nil
}

func (r *SubRepository) UpdateSub(ctx context.Context, id uuid.UUID, req *domain.UpdateSubRequest) (*domain.Sub, error) {
	conn, err := r.pg.GetConnection(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection: %w", err)
	}
	defer conn.Release()

	query := `
			UPDATE subscriptions
			SET
				service_name = COALESCE($1, service_name),
				price = COALESCE($2, price),
				start_date = COALESCE($3, start_date),
				end_date = $4
			WHERE id = $5
			RETURNING id, service_name, price, user_id, start_date, end_date
		`

	var sub domain.Sub
	err = conn.QueryRow(ctx, query,
		req.ServiceName,
		req.Price,
		req.StartDate,
		req.EndDate,
		id,
	).Scan(
		&sub.ID,
		&sub.ServiceName,
		&sub.Price,
		&sub.UserID,
		&sub.StartDate,
		&sub.EndDate,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("error not found: %w", err)
		}
		return nil, fmt.Errorf("failed to update subscription: %w", err)
	}

	return &sub, nil
}

func (r *SubRepository) DeleteSub(ctx context.Context, id uuid.UUID) error {
	conn, err := r.pg.GetConnection(ctx)
	if err != nil {
		return fmt.Errorf("failed to get connection: %w", err)
	}
	defer conn.Release()

	query := `DELETE FROM subscriptions WHERE id = $1`

	cmd, err := conn.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}

	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("error not found: %w", err)

	}

	return nil
}

func (r *SubRepository) CalculateTotalCost(ctx context.Context, filter domain.TotalCostFilter) (int, error) {
	conn, err := r.pg.GetConnection(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get connection: %w", err)
	}
	defer conn.Release()

	query := `
			SELECT COALESCE(SUM(price), 0)
			FROM subscriptions
			WHERE start_date <= $1
			AND (end_date >= $2 OR end_date IS NULL)
		`

	args := []interface{}{filter.EndPeriod, filter.StartPeriod}
	argPos := 3

	if filter.UserID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", argPos)
		args = append(args, *filter.UserID)
		argPos++
	}
	if filter.ServiceName != nil {
		query += fmt.Sprintf(" AND service_name = $%d", argPos)
		args = append(args, *filter.ServiceName)
	}

	var total int
	err = conn.QueryRow(ctx, query, args...).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate total cost: %w", err)
	}

	return total, nil
}
