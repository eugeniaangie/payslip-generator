package service

import (
	"context"

	"payslip-generator/store/postgres_store"
	"payslip-generator/util/errs"

	"github.com/sirupsen/logrus"
)

type PingParams struct {
}

type PingResult struct {
	StatusCode    string `json:"status_code"`
	StatusMessage string `json:"status_msg"`

	Data struct {
		Message string `json:"message"`
	} `json:"data"`
}

func (service *Service) Ping(ctx context.Context, params *PingParams) (*PingResult, error) {
	const op errs.Op = "service/Ping"

	// init service result
	serviceResult := &PingResult{}

	// log the params for data tracing purpose
	service.logger.WithFields(logrus.Fields{
		"op":     op,
		"params": params,
	}).Debug("params!")

	// call data access layer
	storeParams := &postgres_store.PingParam{}

	storeResult, err := service.store.postgres.Ping(ctx, storeParams)

	if err != nil {
		serviceResult.StatusCode = errs.CODE_ERR_IO
		serviceResult.StatusMessage = "nok"

		return serviceResult, err
	}

	// tidy up service result
	serviceResult.StatusCode = errs.CODE_SUCCESS
	serviceResult.StatusMessage = "ok"
	serviceResult.Data.Message = storeResult.Message

	return serviceResult, nil
}
