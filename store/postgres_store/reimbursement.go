package postgres_store

import (
	"context"
	"payslip-generator/model"
	"payslip-generator/util/errs"

	"github.com/sirupsen/logrus"
)

type CreateReimbursementParam struct {
	UserID      string
	Date        string
	Amount      int
	Description string
	CreatedBy   string
	UpdatedBy   string
	IpAddress   string
}

type CreateReimbursementResult struct {
	Data model.Reimbursement
}

func (store *Store) CreateReimbursement(ctx context.Context, param *CreateReimbursementParam) (*CreateReimbursementResult, error) {
	const op errs.Op = "postgres_store/CreateReimbursement"

	query := `
		INSERT INTO reimbursement (user_id, date, amount, description, created_by, updated_by, ip_address)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, user_id, date, amount, description, created_by, ip_address, created_at, updated_by, updated_at
	`

	var reimbursement model.Reimbursement

	err := store.Db.QueryRowxContext(ctx, query,
		param.UserID,
		param.Date,
		param.Amount,
		param.Description,
		param.CreatedBy,
		param.UpdatedBy,
		param.IpAddress,
	).StructScan(&reimbursement)

	if err != nil {
		store.logger.WithFields(logrus.Fields{
			"op":    op,
			"scope": "QueryRowxContext",
			"err":   err.Error(),
		}).Error()
		return nil, err
	}

	result := &CreateReimbursementResult{
		Data: reimbursement,
	}

	return result, nil
}
