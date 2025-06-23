package postgres_store

import (
	"context"
	"fmt"
	"payslip-generator/model"
	"payslip-generator/util/errs"
	"time"

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

type GetPayslipByUserIDAndPeriodIDParam struct {
	UserID   string
	PeriodID string
}

type GetPayslipByUserIDAndPeriodIDResult struct {
	Data model.Payslip
}

func (store *Store) GetPayslipByUserIDAndPeriodID(ctx context.Context, param *GetPayslipByUserIDAndPeriodIDParam) (*GetPayslipByUserIDAndPeriodIDResult, error) {
	const op errs.Op = "postgres_store/GetPayslipByUserIDAndPeriodID"

	query := `
		SELECT * FROM payslip WHERE user_id = $1 AND period_id = $2
	`

	var payslip model.Payslip

	err := store.Db.GetContext(ctx, &payslip, query, param.UserID, param.PeriodID)
	if err != nil {
		store.logger.WithFields(logrus.Fields{
			"op":    op,
			"scope": "GetContext",
			"err":   err.Error(),
		}).Error()
		return nil, err
	}

	result := &GetPayslipByUserIDAndPeriodIDResult{
		Data: payslip,
	}

	return result, nil
}

type UpdatePayslipParam struct {
	BaseSalary    int
	PresentDays   int
	WorkingDays   int
	OvertimeHours int
	Reimbursement int
	TakeHomePay   int
	UpdatedBy     string
	UpdatedAt     time.Time
	IpAddress     string
	ID            string
}

type UpdatePayslipResult struct {
	Data model.Payslip
}

func (store *Store) UpdatePayslip(ctx context.Context, param *UpdatePayslipParam) (*UpdatePayslipResult, error) {
	const op errs.Op = "postgres_store/UpdatePayslip"

	query := `
		UPDATE 
		payslip 
		SET base_salary = $1, 
		present_days = $2, 
		working_days = $3, 
		overtime_hours = $4, 
		reimbursement = $5, 
		take_home_pay = $6, 
		updated_by = $7, 
		updated_at = $8, 
		ip_address = $9 
		WHERE id = $10
		RETURNING id, user_id, period_id, base_salary, present_days, working_days, overtime_hours, reimbursement, take_home_pay, created_by, updated_by, ip_address, created_at, updated_at
	`

	var payslip model.Payslip

	err := store.Db.QueryRowxContext(ctx, query,
		param.BaseSalary,
		param.PresentDays,
		param.WorkingDays,
		param.OvertimeHours,
		param.Reimbursement,
		param.TakeHomePay,
		param.UpdatedBy,
		param.UpdatedAt,
		param.IpAddress,
		param.ID,
	).StructScan(&payslip)

	if err != nil {
		store.logger.WithFields(logrus.Fields{
			"op":    op,
			"scope": "QueryRowxContext",
			"err":   err.Error(),
		}).Error()
		return nil, err
	}

	result := &UpdatePayslipResult{
		Data: payslip,
	}

	return result, nil
}

type GetPayslipSummaryParam struct {
	StartDate string
	EndDate   string
	Page      int64
	PageSize  int64
}

type GetPayslipSummaryResult struct {
	Data       model.PayslipSummary
	Pagination model.Pagination
}

func (store *Store) GetPayslipSummary(ctx context.Context, param *GetPayslipSummaryParam) (*GetPayslipSummaryResult, error) {
	const op errs.Op = "postgres_store/GetPayslipSummary"

	query := `
		SELECT 
			attendance_period.start_date,
			attendance_period.end_date,
			SUM(payslip.take_home_pay) as total_take_home_pay,
			payslip.user_id as employee_id,
			u.full_name as employee_name,
			payslip.present_days,
			payslip.working_days,
			payslip.overtime_hours,
			payslip.reimbursement,
			payslip.take_home_pay
		FROM payslip 
		inner join attendance_period on payslip.period_id = attendance_period.id 
		inner join users u on payslip.user_id = u.id
		WHERE attendance_period.start_date >= $1 AND attendance_period.end_date <= $2
		GROUP BY attendance_period.start_date, attendance_period.end_date, payslip.user_id, u.full_name,payslip.present_days,
			payslip.working_days,
			payslip.overtime_hours,
			payslip.reimbursement,
			payslip.take_home_pay
	`

	queryParams := []interface{}{param.StartDate, param.EndDate}
	paramIdx := 3

	// Count total rows
	countQuery := fmt.Sprintf(`SELECT COUNT(1) FROM (%s) AS x`, query)

	var totalRows int64
	if err := store.Db.GetContext(ctx, &totalRows, countQuery, queryParams...); err != nil {
		store.logger.WithFields(logrus.Fields{
			"op":    op,
			"scope": "CountQuery",
			"err":   err.Error(),
		}).Error()
		return nil, err
	}

	var totalPages, offset, endRow int64
	if param.Page > 0 && param.PageSize > 0 {
		offset = (param.Page - 1) * param.PageSize
		totalPages = (totalRows + param.PageSize - 1) / param.PageSize
		endRow = offset + param.PageSize
		if endRow > totalRows {
			endRow = totalRows
		}

		// Append LIMIT & OFFSET to query
		query += fmt.Sprintf(" ORDER BY attendance_period.start_date DESC LIMIT $%d OFFSET $%d", paramIdx, paramIdx+1)
		queryParams = append(queryParams, param.PageSize, offset)
	} else {
		query += " ORDER BY start_date DESC"
	}

	var items []model.PayslipSummaryItem

	if err := store.Db.SelectContext(ctx, &items, query, queryParams...); err != nil {
		store.logger.WithFields(logrus.Fields{
			"op":    op,
			"scope": "Select",
			"err":   err.Error(),
		}).Error()
		return nil, err
	}

	hasNext := param.Page < totalPages

	resultData := model.PayslipSummary{
		Items:            items,
		PeriodStartDate:  items[0].PeriodStartDate,
		PeriodEndDate:    items[0].PeriodEndDate,
		TotalTakeHomePay: items[0].TotalTakeHomePay,
	}

	result := &GetPayslipSummaryResult{
		Data: resultData,
		Pagination: model.Pagination{
			TotalPage:   totalPages,
			RowPerPage:  param.PageSize,
			TotalRow:    totalRows,
			StartRow:    offset + 1,
			EndRow:      endRow,
			CurrentPage: param.Page,
			HasNext:     hasNext,
		},
	}

	return result, nil
}
