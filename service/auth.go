package service

import (
	"context"
	"errors"
	"fmt"

	"payslip-generator/model"
	"payslip-generator/store/postgres_store"
	"payslip-generator/util/auth"
	"payslip-generator/util/errs"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type RegisterParam struct {
	FullName string `json:"full_name"`
	Username string `json:"username"`
	Password string `json:"password"`
	UserRole string `json:"user_role"`
}

// Register creates a new user.
func (service *Service) Register(ctx context.Context, param *RegisterParam) (*CreateUserResult, error) {
	const op errs.Op = "service/Register"

	service.logger.WithFields(logrus.Fields{
		"op":    op,
		"param": fmt.Sprintf("%+v", param),
	}).Debug()

	serviceResult := &CreateUserResult{}

	usernameExist, err := service.GetUserByUsername(ctx, param.Username)
	if err != nil {
		service.logger.WithFields(logrus.Fields{
			"op":    op,
			"scope": "GetUserByUsername",
			"err":   err.Error(),
		}).Error()

		serviceResult.StatusCode = errs.CODE_ERR_DATABASE
		serviceResult.StatusMessage = "failed to check username"
		return serviceResult, err
	}

	user, ok := usernameExist.Data.(model.User)
	if ok && user.ID != "" {
		// user exists
		serviceResult.StatusCode = errs.CODE_ERR_VALIDATION
		serviceResult.StatusMessage = "username already exists"
		return serviceResult, errors.New("username already exists")
	}

	// Hash password
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(param.Password), bcrypt.DefaultCost)
	if err != nil {
		serviceResult.StatusCode = errs.CODE_ERR_DATABASE
		serviceResult.StatusMessage = "failed to hash password"
		return serviceResult, err
	}

	storeResult, err := service.store.postgres.CreateUser(ctx, &postgres_store.CreateUserParam{
		Username:     param.Username,
		PasswordHash: string(hashedPwd),
		FullName:     param.FullName,
		UserRole:     "employee",
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

type LoginParam struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Login authenticates the user and returns a JWT token.
func (service *Service) Login(ctx context.Context, param *LoginParam) (string, error) {
	const op errs.Op = "service/Login"

	service.logger.WithFields(logrus.Fields{
		"op":    op,
		"param": fmt.Sprintf("%+v", param),
	}).Debug()

	user, err := service.store.postgres.GetUserByUsername(ctx, param.Username)
	if err != nil {
		return "", errors.New("user not found")
	}

	// Compare passwords
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(param.Password)); err != nil {
		return "", errors.New("invalid password")
	}

	// Generate JWT
	return auth.GenerateJWT(user.ID, user.Username)
}
