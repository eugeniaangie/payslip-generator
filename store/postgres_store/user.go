package postgres_store

import (
	"context"
	"fmt"

	"payslip-generator/model"
	"payslip-generator/util/errs"

	"github.com/sirupsen/logrus"
)

func (store *Store) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	const op errs.Op = "postgres_store/GetUserByUsername"

	var user model.User

	err := store.Db.GetContext(ctx, &user, `SELECT * FROM users WHERE username = $1`, username)
	if err != nil {
		store.logger.WithFields(logrus.Fields{
			"op":    op,
			"scope": "QueryRowxContext",
			"err":   err.Error(),
		}).Error()
		return nil, err
	}

	return &user, nil
}

type CreateUserParam struct {
	Username     string
	PasswordHash string
	FullName     string
	UserRole     string
}

type CreateUserResult struct {
	Data model.User
}

func (store *Store) CreateUser(ctx context.Context, param *CreateUserParam) (*CreateUserResult, error) {
	const op errs.Op = "postgres_store/CreateUser"

	query := `
		INSERT INTO users (username, password_hash, full_name, user_role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, username, full_name, user_role, created_at, updated_at
	`

	var user model.User

	err := store.Db.QueryRowxContext(ctx, query,
		param.Username,
		param.PasswordHash,
		param.FullName,
		param.UserRole,
	).StructScan(&user)

	if err != nil {
		store.logger.WithFields(logrus.Fields{
			"op":    op,
			"scope": "QueryRowxContext",
			"err":   err.Error(),
		}).Error()
		return nil, err
	}

	result := &CreateUserResult{
		Data: user,
	}

	return result, nil
}

type GetUserListParam struct {
	Search   string
	Page     int64
	PageSize int64
}

type GetUserListResult struct {
	Data       []model.User
	Pagination model.Pagination
}

func (store *Store) GetUserList(ctx context.Context, param *GetUserListParam) (*GetUserListResult, error) {
	const op errs.Op = "postgres_store/GetUser"

	query := `
	SELECT * FROM users
	WHERE 1=1
	`
	queryParams := []interface{}{}
	paramIdx := 1

	// Filtering by search
	if param.Search != "" {
		query += fmt.Sprintf(" AND full_name ILIKE $%d", paramIdx)
		queryParams = append(queryParams, fmt.Sprintf("%%%s%%", param.Search))
		paramIdx++
	}

	// Count total rows
	countQuery := fmt.Sprintf(`SELECT COUNT(1) FROM (%s) AS x`, query)

	var totalRows int64
	if err := store.Db.GetContext(ctx, &totalRows, countQuery, queryParams...); err != nil {
		store.logger.WithFields(logrus.Fields{
			"op":    op,
			"scope": "CountQuery",
			"err":   err.Error(),
		}).Error()
		return nil, err
	}

	var totalPages, offset, endRow int64
	if param.Page > 0 && param.PageSize > 0 {
		offset = (param.Page - 1) * param.PageSize
		totalPages = (totalRows + param.PageSize - 1) / param.PageSize
		endRow = offset + param.PageSize
		if endRow > totalRows {
			endRow = totalRows
		}

		// Append LIMIT & OFFSET to query
		query += fmt.Sprintf(" ORDER BY full_name ASC LIMIT $%d OFFSET $%d", paramIdx, paramIdx+1)
		queryParams = append(queryParams, param.PageSize, offset)
	} else {
		query += " ORDER BY full_name ASC"
	}

	// Fetch result
	items := []model.User{}
	if err := store.Db.SelectContext(ctx, &items, query, queryParams...); err != nil {
		store.logger.WithFields(logrus.Fields{
			"op":    op,
			"scope": "Select",
			"err":   err.Error(),
		}).Error()
		return nil, err
	}

	hasNext := param.Page < totalPages

	result := &GetUserListResult{
		Data: items,
		Pagination: model.Pagination{
			TotalPage:   totalPages,
			RowPerPage:  param.PageSize,
			TotalRow:    totalRows,
			StartRow:    offset + 1,
			EndRow:      endRow,
			CurrentPage: param.Page,
			HasNext:     hasNext,
		},
	}

	return result, nil
}
