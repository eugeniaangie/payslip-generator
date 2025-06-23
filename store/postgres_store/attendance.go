package postgres_store

import (
	"context"
	"payslip-generator/model"
	"payslip-generator/util/errs"

	"github.com/sirupsen/logrus"
)

type CreateAttendanceParam struct {
	UserID    string
	Date      string
	IpAddress string
	CreatedBy string
	UpdatedBy string
}

type CreateAttendanceResult struct {
	Data model.Attendance
}

func (store *Store) CreateAttendance(ctx context.Context, param *CreateAttendanceParam) (*CreateAttendanceResult, error) {
	const op errs.Op = "postgres_store/CreateAttendance"

	query := `
		INSERT INTO attendance (user_id, date, ip_address, created_by, updated_by)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, date, ip_address, created_by, created_at, updated_by, updated_at
	`

	var attendance model.Attendance

	err := store.Db.QueryRowxContext(ctx, query,
		param.UserID,
		param.Date,
		param.IpAddress,
		param.CreatedBy,
		param.CreatedBy,
	).StructScan(&attendance)

	if err != nil {
		store.logger.WithFields(logrus.Fields{
			"op":    op,
			"scope": "QueryRowxContext",
			"err":   err.Error(),
		}).Error()
		return nil, err
	}

	result := &CreateAttendanceResult{
		Data: attendance,
	}

	return result, nil
}

func (store *Store) GetAttendanceByUserAndDate(ctx context.Context, userID string, date string) (*model.Attendance, error) {
	const op errs.Op = "postgres_store/GetAttendanceByUserAndDate"

	query := `
		SELECT id, user_id, date, ip_address, created_by, created_at, updated_by, updated_at
		FROM attendance
		WHERE user_id = $1 AND date = $2
	`

	var attendance model.Attendance

	err := store.Db.QueryRowxContext(ctx, query, userID, date).StructScan(&attendance)
	if err != nil {
		store.logger.WithFields(logrus.Fields{
			"op":    op,
			"scope": "QueryRowxContext",
			"err":   err.Error(),
		}).Error()
		return nil, err
	}

	return &attendance, nil
}
