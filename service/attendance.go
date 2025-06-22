package service

import (
	"context"
	"fmt"
	"payslip-generator/store/postgres_store"
	"payslip-generator/util/errs"
	"time"

	"github.com/sirupsen/logrus"
)

type CreateAttendanceParam struct {
	UserID    string `json:"user_id"`
	IpAddress string `json:"ip_address"`
	CreatedBy string `json:"created_by"`
}

type CreateAttendanceResult struct {
	StatusCode    string      `json:"status_code"`
	StatusMessage string      `json:"status_msg"`
	Data          interface{} `json:"data"`
}

func (service *Service) CreateAttendance(ctx context.Context, param *CreateAttendanceParam) (*CreateAttendanceResult, error) {
	const op errs.Op = "service/CreateAttendance"

	service.logger.WithFields(logrus.Fields{
		"op":    op,
		"param": fmt.Sprintf("%+v", param),
	}).Debug()

	serviceResult := &CreateAttendanceResult{}

	now := time.Now().In(time.FixedZone("Asia/Jakarta", 7*3600))
	if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
		serviceResult.StatusCode = errs.CODE_SUCCESS
		serviceResult.StatusMessage = "cannot submit on weekends"
		serviceResult.Data = map[string]interface{}{}
		return serviceResult, nil
	}

	dateStr := now.Format("2006-01-02")

	existing, err := service.store.postgres.GetAttendanceByUserAndDate(ctx, param.UserID, dateStr)
	if err != nil && err.Error() != "sql: no rows in result set" {
		serviceResult.StatusCode = errs.CODE_ERR_DATABASE
		serviceResult.StatusMessage = "failed to check attendance"
		return serviceResult, err
	}

	if existing != nil {
		serviceResult.StatusCode = errs.CODE_SUCCESS
		serviceResult.StatusMessage = "attendance already exists"
		serviceResult.Data = *existing
		return serviceResult, nil
	}

	storeResult, err := service.store.postgres.CreateAttendance(ctx, &postgres_store.CreateAttendanceParam{
		UserID:    param.UserID,
		Date:      dateStr,
		IpAddress: param.IpAddress,
		CreatedBy: param.CreatedBy,
	})
	if err != nil {
		service.logger.WithFields(logrus.Fields{
			"op":    op,
			"scope": "CreateAttendance",
			"err":   err.Error(),
		}).Error()

		serviceResult.StatusCode = errs.CODE_ERR_DATABASE
		serviceResult.StatusMessage = "failed to create attendance"
		return serviceResult, err
	}

	serviceResult.StatusCode = errs.CODE_SUCCESS
	serviceResult.StatusMessage = "ok"
	serviceResult.Data = storeResult.Data

	return serviceResult, nil
}
