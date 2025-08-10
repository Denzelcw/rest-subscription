package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"task_manager/internal/config"
	"task_manager/internal/domain"
	"task_manager/internal/http_server/dto"
	"task_manager/internal/storage"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq" // init postgres driver
)

type Storage struct {
	DB *sql.DB
}

func New(dbConfig config.DbConfig) (*Storage, error) {
	const op = "storage.postgresql.New"

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.Host, dbConfig.Port, dbConfig.Username, dbConfig.Password, dbConfig.DBName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{DB: db}, nil
}

func (s *Storage) AddUserSubscription(ctx context.Context, dto dto.CreateUserSubDTO) (int64, error) {
	const op = "storage.postgres.AddUserSubscription"

	const query = `
		INSERT INTO user_subscriptions (
			service_name,
			price,
			user_id,
			start_date,
			end_date
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	startDate, err := time.Parse("01-2006", dto.StartDate)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	var endDatePtr *time.Time
	if dto.EndDate != "" {
		endDate, err := time.Parse("01-2006", dto.EndDate)
		if err != nil {
			return 0, fmt.Errorf("%s: %w", op, err)
		}
		endDatePtr = &endDate
	}

	var id int64
	err = s.DB.QueryRowContext(
		ctx,
		query,
		dto.ServiceName,
		dto.Price,
		dto.UserID,
		startDate,
		endDatePtr,
	).Scan(&id)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUserSubExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetUserSubscriptionById(ctx context.Context, id int) (*domain.UserSubscription, error) {
	const op = "storage.postgresql.GetUserSubscriptionById"

	const query = `
		SELECT 
			id,
			service_name,
			price,
			user_id,
			TO_CHAR(start_date, 'MM-YYYY') AS start_date,
			TO_CHAR(end_date, 'MM-YYYY')   AS end_date
		FROM user_subscriptions
		WHERE id = $1
	`

	var sub domain.UserSubscription
	var endDate sql.NullString

	err := s.DB.QueryRowContext(ctx, query, id).Scan(
		&sub.ID,
		&sub.ServiceName,
		&sub.Price,
		&sub.UserID,
		&sub.StartDate,
		&endDate,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if endDate.Valid {
		sub.EndDate = endDate.String
	}

	return &sub, nil
}

func (s *Storage) GetUserSubscriptionsListByUUID(ctx context.Context, userID uuid.UUID) ([]*domain.UserSubscription, error) {
	const op = "storage.postgresql.GetUserSubscriptionsListByUUID"

	const query = `
		SELECT
			id,
			service_name,
			price,
			user_id,
			TO_CHAR(start_date, 'MM-YYYY') AS start_date,
			TO_CHAR(end_date, 'MM-YYYY') AS end_date
		FROM user_subscriptions
		WHERE user_id = $1
	`

	rows, err := s.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var subscriptions []*domain.UserSubscription

	for rows.Next() {
		var sub domain.UserSubscription
		var endDate sql.NullString

		if err := rows.Scan(
			&sub.ID,
			&sub.ServiceName,
			&sub.Price,
			&sub.UserID,
			&sub.StartDate,
			&endDate,
		); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		if endDate.Valid {
			sub.EndDate = endDate.String
		} else {
			sub.EndDate = ""
		}

		subscriptions = append(subscriptions, &sub)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if len(subscriptions) == 0 {
		return nil, storage.ErrUserNotFound
	}

	return subscriptions, nil
}

func (s *Storage) DeleteUserSubscriptionByID(ctx context.Context, id int) error {
	const op = "storage.postgresql.DeleteUserSubscriptionByID"

	stmt, err := s.DB.Prepare("DELETE FROM user_subscriptions WHERE id = $1")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	result, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected == 0 {
		return storage.ErrNotFound
	}

	return nil
}

func (s *Storage) UpdateUserSubscription(ctx context.Context, dto dto.UpdateUserSubDTO) (*domain.UserSubscription, error) {
	const op = "storage.postgres.UpdateUserSubscription"

	const query = `
		UPDATE user_subscriptions
		SET
			service_name = $2,
			price = $3,
			user_id = $4,
			start_date = $5,
			end_date = $6,
			updated_at = NOW()
		WHERE id = $1
		RETURNING
			id,
			service_name,
			price,
			user_id,
			TO_CHAR(start_date, 'MM-YYYY') AS start_date,
			TO_CHAR(end_date, 'MM-YYYY') AS end_date
	`

	startDate, endDate, err := parseDates(dto.StartDate, dto.EndDate, op)
	if err != nil {
		return nil, err
	}

	var sub domain.UserSubscription

	err = s.DB.QueryRowContext(
		ctx,
		query,
		dto.ID,
		dto.ServiceName,
		dto.Price,
		dto.UserID,
		startDate,
		endDate,
	).Scan(
		&sub.ID,
		&sub.ServiceName,
		&sub.Price,
		&sub.UserID,
		&sub.StartDate,
		&sub.EndDate,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &sub, nil
}

func (s *Storage) CalculateTotalCost(ctx context.Context, dto dto.TotalCost) (int64, error) {
	const op = "storage.postgres.CalculateTotalCost"

	const query = `
		SELECT COALESCE(SUM(price), 0)
		FROM user_subscriptions
		WHERE user_id = $1
		  AND service_name = $2
		  AND start_date >= $3
		  AND end_date <= $4
	`

	startDate, endDate, err := parseDates(dto.StartDate, dto.EndDate, op)
	if err != nil {
		return 0, err
	}

	var totalCost int64
	err = s.DB.QueryRowContext(
		ctx,
		query,
		dto.UserID,
		dto.ServiceName,
		startDate,
		endDate,
	).Scan(&totalCost)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return totalCost, nil
}

func parseDates(startDateStr, endDateStr string, op string) (time.Time, *time.Time, error) {
	var (
		startDate  time.Time
		endDate    time.Time
		endDatePtr *time.Time
		err        error
	)

	startDate, err = time.Parse("01-2006", startDateStr)
	if err != nil {
		return time.Time{}, nil, fmt.Errorf("%s: %w", op, err)
	}

	if endDateStr != "" {
		endDate, err = time.Parse("01-2006", endDateStr)
		if err != nil {
			return time.Time{}, nil, fmt.Errorf("%s: %w", op, err)
		}
		endDatePtr = &endDate
	}

	return startDate, endDatePtr, nil
}
