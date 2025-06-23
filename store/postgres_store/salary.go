package postgres_store

import (
	"context"
	"payslip-generator/model"
	"payslip-generator/util/errs"

	"github.com/sirupsen/logrus"
)

type CreateSalaryParam struct {
	UserID    string
	Amount    int
	IsActive  bool
	CreatedBy string
	UpdatedBy string
}

type CreateSalaryResult struct {
	Data model.Salary
}

func (store *Store) CreateSalary(ctx context.Context, param *CreateSalaryParam) (*CreateSalaryResult, error) {
	const op errs.Op = "postgres_store/CreateSalary"

	query := `
		INSERT INTO salary (user_id, amount, is_active, created_by, updated_by)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, amount, is_active, created_by, updated_by, created_at, updated_at
	`

	var salary model.Salary

	err := store.Db.QueryRowxContext(ctx, query,
		param.UserID,
		param.Amount,
		param.IsActive,
		param.CreatedBy,
		param.UpdatedBy,
	).StructScan(&salary)

	if err != nil {
		store.logger.WithFields(logrus.Fields{
			"op":    op,
			"scope": "QueryRowxContext",
			"err":   err.Error(),
		}).Error()
		return nil, err
	}

	result := &CreateSalaryResult{
		Data: salary,
	}

	return result, nil
}
