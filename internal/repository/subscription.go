package repository

import (
	"context"
	"fmt"
	"testTask/internal/models"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SubscriptionRepository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (s *SubscriptionRepository) Create(ctx context.Context, createModel models.CreateSubscriptionInput) (int, error) {
	startDate, err := time.Parse("01-2006", createModel.StartDate)
	if err != nil {
		return 0, fmt.Errorf("invalid start_date fomat: %v", err)
	}

	var endDate *time.Time
	if createModel.EndDate != nil {
		t, err := time.Parse("01-2006", *createModel.EndDate)
		if err != nil {
			return 0, fmt.Errorf("invalid end_date format: %v", err)
		}
		endDate = &t
	}

	query := `INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id
	`
	var id int

	err = s.db.QueryRow(ctx, query, createModel.ServiceName,
		createModel.Price,
		createModel.UserId,
		startDate,
		endDate,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("unable to make create operation: %v", err)
	}

	return id, nil

}

func (s *SubscriptionRepository) GetByID(ctx context.Context, id int) (*models.Subscription, error) {
	query := `SELECT id, service_name, price, user_id, start_date, end_date
	FROM subscriptions
	WHERE id = $1
	`

	var data models.Subscription

	err := s.db.QueryRow(ctx, query, id).Scan(&data.ID,
		&data.ServiceName,
		&data.Price,
		&data.UserId,
		&data.StartDate,
		&data.EndDate)

	if err != nil {
		return nil, fmt.Errorf("unable to get by id: %v", err)
	}

	return &data, nil

}

func (s *SubscriptionRepository) Update(ctx context.Context, id int, input models.CreateSubscriptionInput) error {
	startDate, err := time.Parse("01-2006", input.StartDate)
	if err != nil {
		return fmt.Errorf("invalid start_Bcdate fomat: %v", err)
	}

	var endDate *time.Time
	if input.EndDate != nil {
		t, err := time.Parse("01-2006", *input.EndDate)
		if err != nil {
			return fmt.Errorf("invalid end_date format: %v", err)
		}
		endDate = &t
	}

	query := `UPDATE subscriptions 
	SET service_name = $1, price = $2, start_date = $3, end_date = $4
	WHERE id = $5
	`
	_, err = s.db.Exec(ctx, query,
		input.ServiceName,
		input.Price,
		startDate,
		endDate,
		id)

	if err != nil {
		return fmt.Errorf("unable to update data: %v", err)
	}

	return nil
}

func (s *SubscriptionRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM subscriptions WHERE id = $1`

	_, err := s.db.Exec(ctx, query, id)

	if err != nil {
		return fmt.Errorf("unable to delete: %v", err)
	}

	return nil
}

func (s *SubscriptionRepository) List(ctx context.Context) ([]models.Subscription, error) {
	query := `SELECT id, service_name, price, user_id, start_date,
  	end_date FROM subscriptions`

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("unable to read data to list: %v", err)
	}
	defer rows.Close()
	var result []models.Subscription

	for rows.Next() {
		var sub models.Subscription
		err := rows.Scan(
			&sub.ID,
			&sub.ServiceName,
			&sub.Price,
			&sub.UserId,
			&sub.StartDate,
			&sub.EndDate)
		if err != nil {
			return nil, fmt.Errorf("unable to list data: %v", err)
		}

		result = append(result, sub)

	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %v", err)
	}
	return result, nil
}

func (s *SubscriptionRepository) TotalCost(ctx context.Context, userID uuid.UUID, serviceName string, from, to time.Time) (int, error) {
	query := `SELECT COALESCE(SUM(price), 0)
  FROM subscriptions
  WHERE user_id = $1
    AND start_date BETWEEN $2 AND $3`

	args := []any{userID, from, to}

	if serviceName != "" {
		query += " AND service_name = $4"
		args = append(args, serviceName)
	}

	var totalCost int
	err := s.db.QueryRow(ctx, query, args...).Scan(&totalCost)
	if err != nil {
		return 0, fmt.Errorf("unable to calculate total cost: %v", err)
	}
	return totalCost, nil

}
