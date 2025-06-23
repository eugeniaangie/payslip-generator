package service

import (
	"context"
	"fmt"
	"payslip-generator/store/postgres_store"
	"payslip-generator/util/errs"

	"github.com/sirupsen/logrus"
)

type CreateReimbursementParam struct {
	UserID      string `json:"user_id"`
	Date        string `json:"date"`
	Amount      int    `json:"amount"`
	Description string `json:"description"`
	CreatedBy   string `json:"created_by"`
	UpdatedBy   string `json:"updated_by"`
	IpAddress   string `json:"ip_address"`
}

type CreateReimbursementResult struct {
	StatusCode    string      `json:"status_code"`
	StatusMessage string      `json:"status_msg"`
	Data          interface{} `json:"data"`
}

func (service *Service) CreateReimbursement(ctx context.Context, param *CreateReimbursementParam) (*CreateReimbursementResult, error) {
	const op errs.Op = "service/CreateReimbursement"

	service.logger.WithFields(logrus.Fields{
		"op":    op,
		"param": fmt.Sprintf("%+v", param),
	}).Debug()

	serviceResult := &CreateReimbursementResult{}

	storeResult, err := service.store.postgres.CreateReimbursement(ctx, &postgres_store.CreateReimbursementParam{
		UserID:      param.UserID,
		Date:        param.Date,
		Amount:      param.Amount,
		Description: param.Description,
		IpAddress:   param.IpAddress,
		CreatedBy:   param.CreatedBy,
		UpdatedBy:   param.UpdatedBy,
	})
	if err != nil {
		service.logger.WithFields(logrus.Fields{
			"op":    op,
			"scope": "CreateReimbursement",
			"err":   err.Error(),
		}).Error()

		serviceResult.StatusCode = errs.CODE_ERR_DATABASE
		serviceResult.StatusMessage = "failed to create reimbursement"
		return serviceResult, err
	}

	serviceResult.StatusCode = errs.CODE_SUCCESS
	serviceResult.StatusMessage = "ok"
	serviceResult.Data = storeResult.Data

	return serviceResult, nil
}
