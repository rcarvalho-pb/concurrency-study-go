package data

import (
	"context"
	"fmt"
	"time"
)

type Plan struct {
	ID         int
	PlanName   string
	PlanAmount int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (p *Plan) GetAll() ([]*Plan, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
	SELECT
		id, plan_name, plan_amount, created_at, updated_at
	FROM 
		plans
	ORDER BY
		plan_name`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var plans []*Plan
	for rows.Next() {
		var plan Plan
		if err := rows.Scan(
			&plan.ID,
			&plan.PlanName,
			&plan.PlanAmount,
			&plan.CreatedAt,
			&plan.UpdatedAt,
		); err != nil {
			return nil, err
		}
		plans = append(plans, &plan)
	}
	return plans, nil
}

func (p *Plan) GetOne(id int) (*Plan, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
	SELECT
		id, plan_name, plan_amount, created_at, updated_at
	FROM
		plans
	WHERE
		id = $1`
	var plan Plan
	row := db.QueryRowContext(ctx, query, id)
	if err := row.Scan(
		&plan.ID,
		&plan.PlanName,
		&plan.PlanAmount,
		&plan.CreatedAt,
		&plan.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &plan, nil
}

func (p *Plan) SubscribeUserToPlan(user User, plan Plan) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	stmt := `
	INSERT INTO
		user_plans
		(user_id, plan_id, created_at, updated_at)
	VALUES
		($1, $2, $3, $4)`
	if _, err := db.ExecContext(ctx, stmt,
		user.ID,
		plan.ID,
		time.Now(),
		time.Now(),
	); err != nil {
		return err
	}
	return nil
}

func (p *Plan) AmountForDisplay() string {
	amount := float64(p.PlanAmount) / 100.0
	return fmt.Sprintf("%.2f", amount)
}

func (p *Plan) GetUserPlan(userID int) (*Plan, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
	SELECT
		p.id, p.plan_name, p.plan_amount, p.created_at, p.updated_at
	FROM 
		user_plans up
	LEFT JOIN 
		plans p
	ON
		p.id = up.plan_id
	WHERE
		up.user_id = $1`

	var plan Plan
	row := db.QueryRowContext(ctx, query, userID)
	if err := row.Scan(
		&plan.ID,
		&plan.PlanName,
		&plan.PlanAmount,
		&plan.CreatedAt,
		&plan.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &plan, nil
}
