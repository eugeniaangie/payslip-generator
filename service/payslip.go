package service

import (
	"context"
	"fmt"
	"payslip-generator/store/postgres_store"
	"payslip-generator/util/errs"
	"time"

	"github.com/sirupsen/logrus"
)

type CreatePayslipParam struct {
	PeriodID  string `json:"period_id"`
	UserID    string `json:"user_id"`
	IpAddress string `json:"ip_address"`
	CreatedBy string `json:"created_by"`
	UpdatedBy string `json:"updated_by"`
}

type CreatePayslipResult struct {
	StatusCode    string      `json:"status_code"`
	StatusMessage string      `json:"status_msg"`
	Data          interface{} `json:"data"`
}

func (service *Service) CreatePayslip(ctx context.Context, param *CreatePayslipParam) (*CreatePayslipResult, error) {
	const op errs.Op = "service/CreatePayslip"

	service.logger.WithFields(logrus.Fields{
		"op":    op,
		"param": fmt.Sprintf("%+v", param),
	}).Debug()

	serviceResult := &CreatePayslipResult{}

	salary, err := service.store.postgres.GetSalaryByUserID(ctx, &postgres_store.GetSalaryByUserIDParam{
		UserID: param.UserID,
	})
	if err != nil {
		service.logger.WithFields(logrus.Fields{
			"op":    op,
			"scope": "GetSalaryByUserID",
			"err":   err.Error(),
		}).Error()

		serviceResult.StatusCode = errs.CODE_ERR_DATABASE
		serviceResult.StatusMessage = "failed to get salary"
		return serviceResult, err
	}

	attendancePeriod, err := service.store.postgres.GetAttendancePeriodByID(ctx, &postgres_store.GetAttendancePeriodByIDParam{
		ID: param.PeriodID,
	})
	if err != nil {
		service.logger.WithFields(logrus.Fields{
			"op":    op,
			"scope": "GetAttendancePeriodByID",
			"err":   err.Error(),
		}).Error()

		serviceResult.StatusCode = errs.CODE_ERR_DATABASE
		serviceResult.StatusMessage = "failed to get attendance period"
		return serviceResult, err
	}

	attendance, err := service.store.postgres.GetAttendanceByUserIDAndPeriodID(ctx, &postgres_store.GetAttendanceByUserIDAndPeriodIDParam{
		UserID:    param.UserID,
		StartDate: attendancePeriod.Data.StartDate,
		EndDate:   attendancePeriod.Data.EndDate,
	})
	if err != nil {
		service.logger.WithFields(logrus.Fields{
			"op":    op,
			"scope": "GetAttendanceByUserIDAndPeriodID",
			"err":   err.Error(),
		}).Error()

		serviceResult.StatusCode = errs.CODE_ERR_DATABASE
		serviceResult.StatusMessage = "failed to get attendance"
		return serviceResult, err
	}

	overtime, err := service.store.postgres.GetOvertimeByUserIDAndPeriodID(ctx, &postgres_store.GetOvertimeByUserIDAndPeriodIDParam{
		UserID:    param.UserID,
		StartDate: attendancePeriod.Data.StartDate,
		EndDate:   attendancePeriod.Data.EndDate,
	})
	if err != nil {
		service.logger.WithFields(logrus.Fields{
			"op":    op,
			"scope": "GetOvertimeByUserIDAndPeriodID",
			"err":   err.Error(),
		}).Error()

		serviceResult.StatusCode = errs.CODE_ERR_DATABASE
		serviceResult.StatusMessage = "failed to get overtime"
		return serviceResult, err
	}

	reimbursement, err := service.store.postgres.GetReimbursementByUserIDAndPeriodID(ctx, &postgres_store.GetReimbursementByUserIDAndPeriodIDParam{
		UserID:    param.UserID,
		StartDate: attendancePeriod.Data.StartDate,
		EndDate:   attendancePeriod.Data.EndDate,
	})
	if err != nil {
		service.logger.WithFields(logrus.Fields{
			"op":    op,
			"scope": "GetReimbursementByUserIDAndPeriodID",
			"err":   err.Error(),
		}).Error()

		serviceResult.StatusCode = errs.CODE_ERR_DATABASE
		serviceResult.StatusMessage = "failed to get reimbursement"
		return serviceResult, err
	}

	workingDays, err := countWeekdays(service, attendancePeriod.Data.StartDate, attendancePeriod.Data.EndDate)
	if err != nil {
		service.logger.WithFields(logrus.Fields{
			"op":    op,
			"scope": "countWeekdays",
			"err":   err.Error(),
		}).Error()

		serviceResult.StatusCode = errs.CODE_ERR_VALIDATION
		serviceResult.StatusMessage = "failed to count weekdays"
		return serviceResult, err
	}

	attendanceDays := len(attendance.Data)

	overtimeHours := 0
	for _, overtime := range overtime.Data {
		overtimeHours += overtime.Hours
	}

	dailySalary := float64(salary.Data.Amount) / float64(workingDays)
	salaryByAttendance := float64(attendanceDays) * dailySalary

	hourlySalary := dailySalary / 8
	overtimePay := float64(overtimeHours) * hourlySalary * 2

	totalReimbursement := 0
	for _, reimbursement := range reimbursement.Data {
		totalReimbursement += reimbursement.Amount
	}

	takeHomePay := int(salaryByAttendance + overtimePay + float64(totalReimbursement))

	storeResult, err := service.store.postgres.CreatePayslip(ctx, &postgres_store.CreatePayslipParam{
		UserID:        param.UserID,
		PeriodID:      param.PeriodID,
		BaseSalary:    salary.Data.Amount,
		PresentDays:   attendanceDays,
		WorkingDays:   workingDays,
		OvertimeHours: overtimeHours,
		Reimbursement: totalReimbursement,
		TakeHomePay:   takeHomePay,
		CreatedBy:     param.CreatedBy,
		UpdatedBy:     param.UpdatedBy,
		IpAddress:     param.IpAddress,
	})
	if err != nil {
		service.logger.WithFields(logrus.Fields{
			"op":    op,
			"scope": "CreatePayslip",
			"err":   err.Error(),
		}).Error()

		serviceResult.StatusCode = errs.CODE_ERR_DATABASE
		serviceResult.StatusMessage = "failed to create payslip"
		return serviceResult, err
	}

	serviceResult.StatusCode = errs.CODE_SUCCESS
	serviceResult.StatusMessage = "ok"
	serviceResult.Data = storeResult.Data

	return serviceResult, nil
}

// Helper: flexibly parse date string in either '2006-01-02' or RFC3339
func parseDateFlexible(s string) (time.Time, error) {
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return t, nil
	}
	return time.Parse(time.RFC3339, s)
}

// Helper: count weekdays between two date strings (inclusive)
func countWeekdays(service *Service, start, end string) (int, error) {
	const op errs.Op = "service/countWeekdays"

	startDate, err := parseDateFlexible(start)
	if err != nil {
		service.logger.WithFields(logrus.Fields{
			"op":    op,
			"scope": "countWeekdays",
			"err":   err.Error(),
		}).Error()
		return 0, err
	}
	endDate, err := parseDateFlexible(end)
	if err != nil {
		service.logger.WithFields(logrus.Fields{
			"op":    op,
			"scope": "countWeekdays",
			"err":   err.Error(),
		}).Error()
		return 0, err
	}

	weekdays := 0
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		if d.Weekday() != time.Saturday && d.Weekday() != time.Sunday {
			weekdays++
		}
	}
	return weekdays, nil
}
