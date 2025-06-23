package service

import (
	"context"
	"fmt"
	"payslip-generator/store/postgres_store"
	"payslip-generator/util/errs"
	"time"

	"github.com/sirupsen/logrus"
)

type RunPayrollParam struct {
	PeriodID  string `json:"period_id"`
	IpAddress string `json:"ip_address"`
	CreatedBy string `json:"created_by"`
	UpdatedBy string `json:"updated_by"`
}

type RunPayrollResult struct {
	StatusCode    string      `json:"status_code"`
	StatusMessage string      `json:"status_msg"`
	Data          interface{} `json:"data"`
}

func (service *Service) RunPayroll(ctx context.Context, param *RunPayrollParam) (*RunPayrollResult, error) {
	const op errs.Op = "service/RunPayroll"

	service.logger.WithFields(logrus.Fields{
		"op":    op,
		"param": fmt.Sprintf("%+v", param),
	}).Debug()

	serviceResult := &RunPayrollResult{}
	data := []interface{}{}

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

	if attendancePeriod.Data.IsPayrollProcessed {
		serviceResult.StatusCode = errs.CODE_ERR_VALIDATION
		serviceResult.StatusMessage = "attendance period already processed"
		return serviceResult, err
	}

	employees, err := service.store.postgres.GetUserList(ctx, &postgres_store.GetUserListParam{
		Page:     1,
		PageSize: 1000,
		UserRole: "employee",
	})
	if err != nil {
		service.logger.WithFields(logrus.Fields{
			"op":    op,
			"scope": "GetUserList",
			"err":   err.Error(),
		}).Error()

		serviceResult.StatusCode = errs.CODE_ERR_DATABASE
		serviceResult.StatusMessage = "failed to get employees"
		return serviceResult, err
	}

	if len(employees.Data) != 0 {
		for _, employee := range employees.Data {
			payslipResult, err := service.CreatePayslip(ctx, &CreatePayslipParam{
				PeriodID:  param.PeriodID,
				UserID:    employee.ID,
				IpAddress: param.IpAddress,
				CreatedBy: param.CreatedBy,
				UpdatedBy: param.UpdatedBy,
			})
			if err != nil {
				service.logger.WithFields(logrus.Fields{
					"op":    op,
					"scope": "CreatePayslip",
					"err":   err.Error(),
				}).Error()

				serviceResult.StatusCode = errs.CODE_ERR_IO
				serviceResult.StatusMessage = "failed to create payslip"
				return serviceResult, err
			}
			data = append(data, payslipResult.Data)
		}
	}

	_, err = service.store.postgres.UpdatePayrollAttendancePeriod(ctx, &postgres_store.UpdatePayrollAttendancePeriodParam{
		ID:                 param.PeriodID,
		IsPayrollProcessed: true,
		UpdatedBy:          param.UpdatedBy,
		UpdatedAt:          time.Now(),
	})
	if err != nil {
		service.logger.WithFields(logrus.Fields{
			"op":    op,
			"scope": "UpdatePayrollAttendancePeriod",
			"err":   err.Error(),
		}).Error()

		serviceResult.StatusCode = errs.CODE_ERR_DATABASE
		serviceResult.StatusMessage = "failed to update attendance period"
		return serviceResult, err
	}

	serviceResult.StatusCode = errs.CODE_SUCCESS
	serviceResult.StatusMessage = "ok"
	serviceResult.Data = data
	return serviceResult, nil
}
