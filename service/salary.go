package service

import (
	"context"
	"fmt"
	"payslip-generator/store/postgres_store"
	"payslip-generator/util/errs"

	"github.com/sirupsen/logrus"
)

type CreateSalaryParam struct {
	UserID    string `json:"user_id"`
	Amount    int    `json:"amount"`
	IsActive  bool   `json:"is_active"`
	CreatedBy string `json:"created_by"`
	UpdatedBy string `json:"updated_by"`
}

type CreateSalaryResult struct {
	StatusCode    string      `json:"status_code"`
	StatusMessage string      `json:"status_msg"`
	Data          interface{} `json:"data"`
}

func (service *Service) CreateSalary(ctx context.Context, param *CreateSalaryParam) (*CreateSalaryResult, error) {
	const op errs.Op = "service/CreateSalary"

	service.logger.WithFields(logrus.Fields{
		"op":    op,
		"param": fmt.Sprintf("%+v", param),
	}).Debug()

	serviceResult := &CreateSalaryResult{}

	storeResult, err := service.store.postgres.CreateSalary(ctx, &postgres_store.CreateSalaryParam{
		UserID:    param.UserID,
		Amount:    param.Amount,
		IsActive:  param.IsActive,
		CreatedBy: param.CreatedBy,
		UpdatedBy: param.UpdatedBy,
	})
	if err != nil {
		service.logger.WithFields(logrus.Fields{
			"op":    op,
			"scope": "CreateSalary",
			"err":   err.Error(),
		}).Error()

		serviceResult.StatusCode = errs.CODE_ERR_DATABASE
		serviceResult.StatusMessage = "failed to create salary"
		return serviceResult, err
	}

	serviceResult.StatusCode = errs.CODE_SUCCESS
	serviceResult.StatusMessage = "ok"
	serviceResult.Data = storeResult.Data

	return serviceResult, nil
}
