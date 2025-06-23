package postgres_store

import (
	"context"
	"payslip-generator/model"
	"payslip-generator/util/errs"

	"github.com/sirupsen/logrus"
)

type CreatePayslipParam struct {
	UserID        string
	PeriodID      string
	BaseSalary    int
	PresentDays   int
	WorkingDays   int
	OvertimeHours int
	Reimbursement int
	TakeHomePay   int
	CreatedBy     string
	UpdatedBy     string
	IpAddress     string
}

type CreatePayslipResult struct {
	Data model.Payslip
}

func (store *Store) CreatePayslip(ctx context.Context, param *CreatePayslipParam) (*CreatePayslipResult, error) {
	const op errs.Op = "postgres_store/CreatePayslip"

	query := `
		INSERT INTO payslip (user_id, period_id, base_salary, present_days, working_days, overtime_hours, reimbursement, take_home_pay, created_by, updated_by, ip_address)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, user_id, period_id, base_salary, present_days, working_days, overtime_hours, reimbursement, take_home_pay, created_by, updated_by, ip_address, created_at, updated_at
	`

	var payslip model.Payslip

	err := store.Db.QueryRowxContext(ctx, query,
		param.UserID,
		param.PeriodID,
		param.BaseSalary,
		param.PresentDays,
		param.WorkingDays,
		param.OvertimeHours,
		param.Reimbursement,
		param.TakeHomePay,
		param.CreatedBy,
		param.UpdatedBy,
		param.IpAddress,
	).StructScan(&payslip)

	if err != nil {
		store.logger.WithFields(logrus.Fields{
			"op":    op,
			"scope": "QueryRowxContext",
			"err":   err.Error(),
		}).Error()
		return nil, err
	}

	result := &CreatePayslipResult{
		Data: payslip,
	}

	return result, nil
}
