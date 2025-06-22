package service

import (
	"context"
	"fmt"
	"payslip-generator/model"
	"payslip-generator/store/postgres_store"
	"payslip-generator/util/errs"

	"github.com/sirupsen/logrus"
)

type CreateAttendancePeriodParam struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	CreatedBy string `json:"created_by"`
	UpdatedBy string `json:"updated_by"`
}

type CreateAttendancePeriodResult struct {
	StatusCode    string                 `json:"status_code"`
	StatusMessage string                 `json:"status_msg"`
	Data          model.AttendancePeriod `json:"data"`
}

func (service *Service) CreateAttendancePeriod(ctx context.Context, param *CreateAttendancePeriodParam) (*CreateAttendancePeriodResult, error) {
	const op errs.Op = "service/CreateAttendancePeriod"
	
	service.logger.WithFields(logrus.Fields{
		"op":    op,
		"param": fmt.Sprintf("%+v", param),
	}).Debug()

	serviceResult := &CreateAttendancePeriodResult{}

	storeResult, err := service.store.postgres.CreateAttendancePeriod(ctx, &postgres_store.CreateAttendancePeriodParam{
		StartDate: param.StartDate,
		EndDate:   param.EndDate,
		CreatedBy: param.CreatedBy,
		UpdatedBy: param.UpdatedBy,
	})
	if err != nil {
		service.logger.WithFields(logrus.Fields{
			"op":    op,
			"scope": "CreateAttendancePeriod",
			"err":   err.Error(),
		}).Error()

		serviceResult.StatusCode = errs.CODE_ERR_DATABASE
		serviceResult.StatusMessage = "failed to create attendance period"
		return serviceResult, err
	}

	serviceResult.StatusCode = errs.CODE_SUCCESS
	serviceResult.StatusMessage = "ok"
	serviceResult.Data = storeResult.Data

	return serviceResult, nil
}

type GetAttendancePeriodListParam struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Page      int64  `json:"page"`
	PageSize  int64  `json:"page_size"`
}

type GetAttendancePeriodListResult struct {
	StatusCode    string                   `json:"status_code"`
	StatusMessage string                   `json:"status_msg"`
	Data          []model.AttendancePeriod `json:"data"`
	Pagination    model.Pagination         `json:"pagination"`
}

func (service *Service) GetAttendancePeriodList(ctx context.Context, param *GetAttendancePeriodListParam) (*GetAttendancePeriodListResult, error) {
	const op errs.Op = "service/GetAttendancePeriodList"

	serviceResult := &GetAttendancePeriodListResult{}

	if param.Page <= 0 {
		param.Page = model.DefaultPage
	}
	if param.PageSize <= 0 {
		param.PageSize = model.DefaultPageSize
	}
	// Limit maximum page size
	if param.PageSize > model.MaxPageSize {
		param.PageSize = model.MaxPageSize
	}

	service.logger.WithFields(logrus.Fields{
		"op":    op,
		"param": fmt.Sprintf("%+v", param),
	}).Debug()

	storeResult, err := service.store.postgres.GetAttendancePeriodList(ctx, &postgres_store.GetAttendancePeriodListParam{
		Page:      param.Page,
		PageSize:  param.PageSize,
		StartDate: param.StartDate,
		EndDate:   param.EndDate,
	})
	if err != nil {
		service.logger.WithFields(logrus.Fields{
			"op":    op,
			"scope": "GetAttendancePeriodList",
			"err":   err.Error(),
		}).Error()

		serviceResult.StatusCode = errs.CODE_ERR_DATABASE
		serviceResult.StatusMessage = "nok"

		return serviceResult, err
	}

	serviceResult.StatusCode = errs.CODE_SUCCESS
	serviceResult.StatusMessage = "ok"
	serviceResult.Data = storeResult.Data
	serviceResult.Pagination = storeResult.Pagination

	return serviceResult, nil
}
