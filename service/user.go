package service

import (
	"context"
	"fmt"
	"payslip-generator/model"
	"payslip-generator/store/postgres_store"
	"payslip-generator/util/errs"

	"github.com/sirupsen/logrus"
)

type GetUserByUsernameParam struct {
	Username string `json:"username"`
}

type GetUserByUsernameResult struct {
	StatusCode    string      `json:"status_code"`
	StatusMessage string      `json:"status_msg"`
	Data          interface{} `json:"data"`
}

func (service *Service) GetUserByUsername(ctx context.Context, param *GetUserByUsernameParam) (*GetUserByUsernameResult, error) {
	const op errs.Op = "service/GetUserByUsername"

	service.logger.WithFields(logrus.Fields{
		"op":    op,
		"param": fmt.Sprintf("%+v", param),
	}).Debug()

	serviceResult := &GetUserByUsernameResult{}

	storeResult, err := service.store.postgres.GetUserByUsername(ctx, param.Username)
	if err != nil {
		// Check if it's a "no rows" error
		if err.Error() == "sql: no rows in result set" {
			serviceResult.StatusCode = errs.CODE_SUCCESS
			serviceResult.StatusMessage = "user not found"
			serviceResult.Data = map[string]interface{}{} // Empty user
			return serviceResult, nil
		}

		service.logger.WithFields(logrus.Fields{
			"op":    op,
			"scope": "GetUserByUsername",
			"err":   err.Error(),
		}).Error()

		serviceResult.StatusCode = errs.CODE_ERR_DATABASE
		serviceResult.StatusMessage = "failed to get user by username"
		return serviceResult, err
	}

	serviceResult.StatusCode = errs.CODE_SUCCESS
	serviceResult.StatusMessage = "ok"
	serviceResult.Data = *storeResult

	return serviceResult, nil
}

type CreateUserParam struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	FullName     string `json:"full_name"`
	UserRole     string `json:"user_role"`
}

type CreateUserResult struct {
	StatusCode    string     `json:"status_code"`
	StatusMessage string     `json:"status_msg"`
	Data          model.User `json:"data"`
}

func (service *Service) CreateUser(ctx context.Context, param *CreateUserParam) (*CreateUserResult, error) {
	const op errs.Op = "service/CreateUser"

	service.logger.WithFields(logrus.Fields{
		"op":    op,
		"param": fmt.Sprintf("%+v", param),
	}).Debug()

	serviceResult := &CreateUserResult{}

	storeResult, err := service.store.postgres.CreateUser(ctx, &postgres_store.CreateUserParam{
		Username:     param.Username,
		PasswordHash: param.PasswordHash,
		FullName:     param.FullName,
		UserRole:     param.UserRole,
	})
	if err != nil {
		service.logger.WithFields(logrus.Fields{
			"op":    op,
			"scope": "CreateUser",
			"err":   err.Error(),
		}).Error()

		serviceResult.StatusCode = errs.CODE_ERR_DATABASE
		serviceResult.StatusMessage = "failed to create user"
		return serviceResult, err
	}

	serviceResult.StatusCode = errs.CODE_SUCCESS
	serviceResult.StatusMessage = "ok"
	serviceResult.Data = storeResult.Data

	return serviceResult, nil
}

type GetUserListParam struct {
	Search   string `json:"search"`
	Page     int64  `json:"page"`
	PageSize int64  `json:"page_size"`
}

type GetUserListResult struct {
	StatusCode    string           `json:"status_code"`
	StatusMessage string           `json:"status_msg"`
	Data          []model.User     `json:"data"`
	Pagination    model.Pagination `json:"pagination"`
}

func (service *Service) GetUserList(ctx context.Context, param *GetUserListParam) (*GetUserListResult, error) {
	const op errs.Op = "service/GetUserList"

	serviceResult := &GetUserListResult{}

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

	storeResult, err := service.store.postgres.GetUserList(ctx, &postgres_store.GetUserListParam{
		Page:     param.Page,
		PageSize: param.PageSize,
		Search:   param.Search,
	})
	if err != nil {
		service.logger.WithFields(logrus.Fields{
			"op":    op,
			"scope": "GetUserList",
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
