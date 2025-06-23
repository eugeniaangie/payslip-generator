package service

import (
	"context"
	"fmt"
	"payslip-generator/store/postgres_store"
	"payslip-generator/util/errs"
	"time"

	"github.com/sirupsen/logrus"
)

type CreateOvertimeParam struct {
	UserID    string `json:"user_id"`
	Hours     int    `json:"hours"`
	IpAddress string `json:"ip_address"`
	CreatedBy string `json:"created_by"`
	UpdatedBy string `json:"updated_by"`
}

type CreateOvertimeResult struct {
	StatusCode    string      `json:"status_code"`
	StatusMessage string      `json:"status_msg"`
	Data          interface{} `json:"data"`
}

func (service *Service) CreateOvertime(ctx context.Context, param *CreateOvertimeParam) (*CreateOvertimeResult, error) {
	const op errs.Op = "service/CreateOvertime"

	service.logger.WithFields(logrus.Fields{
		"op":    op,
		"param": fmt.Sprintf("%+v", param),
	}).Debug()

	serviceResult := &CreateOvertimeResult{}

	now := time.Now().In(time.FixedZone("Asia/Jakarta", 7*3600))
	dateStr := now.Format("2006-01-02")

	// Validate date not in the future
	if dateStr > now.Format("2006-01-02") {
		serviceResult.StatusCode = errs.CODE_ERR_VALIDATION
		serviceResult.StatusMessage = "date not in the future"
		return serviceResult, nil
	}

	// Validate date with regular 8 working hours per day (9AM-5PM), 5 days a week (monday-friday)
	if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday || now.Hour() < 9 || now.Hour() > 17 {
		serviceResult.StatusCode = errs.CODE_SUCCESS
		serviceResult.StatusMessage = "cannot submit on weekends"
		serviceResult.Data = map[string]interface{}{}
		return serviceResult, nil
	}

	// Validate hours between 1-3
	if param.Hours < 1 || param.Hours > 3 {
		serviceResult.StatusCode = errs.CODE_ERR_VALIDATION
		serviceResult.StatusMessage = "hours must be between 1-3"
		return serviceResult, nil
	}

	// Validate date not duplicate
	existing, err := service.store.postgres.GetOvertimeByUserAndDate(ctx, param.UserID, dateStr)
	if err != nil && err.Error() != "sql: no rows in result set" {
		serviceResult.StatusCode = errs.CODE_ERR_DATABASE
		serviceResult.StatusMessage = "failed to check attendance"
		return serviceResult, err
	}

	if existing != nil {
		serviceResult.StatusCode = errs.CODE_SUCCESS
		serviceResult.StatusMessage = "overtime already exists"
		serviceResult.Data = *existing
		return serviceResult, nil
	}

	storeParams := &postgres_store.CreateOvertimeParam{
		UserID:    param.UserID,
		Date:      dateStr,
		Hours:     param.Hours,
		IpAddress: param.IpAddress,
		CreatedBy: param.CreatedBy,
		UpdatedBy: param.UpdatedBy,
	}

	storeResult, err := service.store.postgres.CreateOvertime(ctx, storeParams)
	if err != nil {
		serviceResult.StatusCode = errs.CODE_ERR_DATABASE
		serviceResult.StatusMessage = "failed to create overtime"
		return serviceResult, err
	}

	serviceResult.StatusCode = errs.CODE_SUCCESS
	serviceResult.StatusMessage = "ok"
	serviceResult.Data = storeResult.Data

	return serviceResult, nil
}
