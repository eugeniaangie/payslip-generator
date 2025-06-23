package postgres_store

import (
	"context"
	"payslip-generator/model"
	"payslip-generator/util/errs"

	"github.com/sirupsen/logrus"
)

type CreateOvertimeParam struct {
	UserID    string
	Date      string
	Hours     int
	IpAddress string
	CreatedBy string
	UpdatedBy string
}

type CreateOvertimeResult struct {
	Data model.Overtime
}

func (store *Store) CreateOvertime(ctx context.Context, param *CreateOvertimeParam) (*CreateOvertimeResult, error) {
	const op errs.Op = "postgres_store/CreateOvertime"

	query := `
		INSERT INTO overtime (user_id, date, hours, ip_address, created_by, updated_by)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, user_id, date, hours, ip_address, created_by, created_at, updated_by, updated_at
	`

	var overtime model.Overtime

	err := store.Db.QueryRowxContext(ctx, query,
		param.UserID,
		param.Date,
		param.Hours,
		param.IpAddress,
		param.CreatedBy,
		param.CreatedBy,
	).StructScan(&overtime)

	if err != nil {
		store.logger.WithFields(logrus.Fields{
			"op":    op,
			"scope": "QueryRowxContext",
			"err":   err.Error(),
		}).Error()
		return nil, err
	}

	result := &CreateOvertimeResult{
		Data: overtime,
	}

	return result, nil
}

func (store *Store) GetOvertimeByUserAndDate(ctx context.Context, userID string, date string) (*model.Overtime, error) {
	const op errs.Op = "postgres_store/GetOvertimeByUserAndDate"

	query := `
		SELECT id, user_id, date, ip_address, created_by, created_at, updated_by, updated_at
		FROM overtime
		WHERE user_id = $1 AND date = $2
	`

	var overtime model.Overtime

	err := store.Db.QueryRowxContext(ctx, query, userID, date).StructScan(&overtime)
	if err != nil {
		store.logger.WithFields(logrus.Fields{
			"op":    op,
			"scope": "QueryRowxContext",
			"err":   err.Error(),
		}).Error()
		return nil, err
	}

	return &overtime, nil
}
