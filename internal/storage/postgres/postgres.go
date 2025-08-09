package postgres

import (
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

func (s *Storage) AddSubscription(dto dto.CreateSubDTO) (int64, error) {
	const op = "storage.postgres.AddSubscription"

	stmt, err := s.DB.Prepare(
		"INSERT INTO user_subscriptions(service_name, price, user_id, start_date, end_date) VALUES($1, $2, $3, $4, $5) RETURNING id;")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	var id int64
	startDate, err := time.Parse("01-2006", dto.StartDate)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	var endDate time.Time
	var endDatePtr *time.Time
	if dto.EndDate != "" {
		endDate, err = time.Parse("01-2006", dto.EndDate)
		if err != nil {
			return 0, fmt.Errorf("%s: %w", op, err)
		}
		endDatePtr = &endDate
	}

	err = stmt.QueryRow(dto.ServiceName, dto.Price, dto.UserID, startDate, endDatePtr).Scan(&id)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return 0, fmt.Errorf("%s: %w", op, storage.ErrUrlExists)
			}
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetSubscriptionById(id int) (*domain.Subscription, error) {
	const op = "storage.postgresql.GetSubscriptionById"

	stmt, err := s.DB.Prepare("SELECT id, service_name, price, user_id, TO_CHAR(start_date, 'MM-YYYY') AS start_date, TO_CHAR(end_date, 'MM-YYYY') AS end_date FROM user_subscriptions WHERE id = $1")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var sub domain.Subscription
	var endDate sql.NullString

	err = stmt.QueryRow(id).Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &endDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if endDate.Valid {
		sub.EndDate = endDate.String
	} else {
		sub.EndDate = ""
	}

	return &sub, nil
}

func (s *Storage) GetSubscriptionsListByUUID(userID uuid.UUID) ([]*domain.Subscription, error) {
	const op = "storage.postgresql.GetSubscriptionsListByUUID"

	stmt, err := s.DB.Prepare(`
  SELECT
   id,
   service_name,
   price,
   user_id,
   TO_CHAR(start_date, 'MM-YYYY') AS start_date,
   TO_CHAR(end_date, 'MM-YYYY') AS end_date
  FROM
   user_subscriptions
  WHERE
   user_id = $1
 `)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := stmt.Query(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var subscriptions []*domain.Subscription

	for rows.Next() {
		var sub domain.Subscription
		var endDate sql.NullString

		err = rows.Scan(
			&sub.ID,
			&sub.ServiceName,
			&sub.Price,
			&sub.UserID,
			&sub.StartDate,
			&endDate,
		)
		if err != nil {
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

func (s *Storage) DeleteSubscriptionByID(id int) error {
	const op = "storage.postgresql.DeleteSubscriptionByID"

	stmt, err := s.DB.Prepare("DELETE FROM user_subscriptions WHERE id = $1")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	result, err := stmt.Exec(id)
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

func (s *Storage) UpdateSubscription(dto dto.UpdateSubDTO) (*domain.Subscription, error) {
	const op = "storage.postgres.UpdateSubscription"

	stmt, err := s.DB.Prepare(`
	  UPDATE user_subscriptions
	  SET service_name = $2,
	   price = $3,
	   user_id = $4,
	   start_date = $5,
	   end_date = $6,
	   updated_at = NOW()
	  WHERE id = $1
	  RETURNING id, service_name, price, user_id, TO_CHAR(start_date, 'MM-YYYY') AS start_date, TO_CHAR(end_date, 'MM-YYYY') AS end_date;
 	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	startDate, endDate, err := parseDates(dto.StartDate, dto.EndDate, op)

	var sub domain.Subscription

	err = stmt.QueryRow(dto.ID, dto.ServiceName, dto.Price, dto.UserID, startDate, endDate).Scan(
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

func (s *Storage) CalculateTotalCost(dto dto.TotalCost) (int64, error) {
	const op = "handler.CalculateTotalCost"
	query := `
		SELECT COALESCE(SUM(price), 0)
		FROM user_subscriptions
		WHERE user_id = $1
		AND service_name = $2
		AND start_date >= $3
		AND end_date <= $4
 	`

	startDate, endDate, err := parseDates(dto.StartDate, dto.EndDate, op)

	row := s.DB.QueryRow(query, dto.UserID, dto.ServiceName, startDate, endDate)

	var totalCost int64

	err = row.Scan(&totalCost)
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
