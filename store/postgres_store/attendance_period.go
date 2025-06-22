package postgres_store

import (
	"context"
	"fmt"

	"payslip-generator/model"
	"payslip-generator/util/errs"

	"github.com/sirupsen/logrus"
)

type CreateAttendancePeriodParam struct {
	StartDate string
	EndDate   string
	CreatedBy string
	UpdatedBy string
}

type CreateAttendancePeriodResult struct {
	Data model.AttendancePeriod
}

func (store *Store) CreateAttendancePeriod(ctx context.Context, param *CreateAttendancePeriodParam) (*CreateAttendancePeriodResult, error) {
	const op errs.Op = "postgres_store/CreateAttendancePeriod"

	query := `
		INSERT INTO attendance_period (start_date, end_date, created_by, updated_by)
		VALUES ($1, $2, $3, $4)
		RETURNING id, start_date, end_date, created_by, updated_by, created_at, updated_at
	`

	var attendancePeriod model.AttendancePeriod

	err := store.Db.QueryRowxContext(ctx, query,
		param.StartDate,
		param.EndDate,
		param.CreatedBy,
		param.UpdatedBy,
	).StructScan(&attendancePeriod)

	if err != nil {
		store.logger.WithFields(logrus.Fields{
			"op":    op,
			"scope": "QueryRowxContext",
			"err":   err.Error(),
		}).Error()
		return nil, err
	}

	result := &CreateAttendancePeriodResult{
		Data: attendancePeriod,
	}

	return result, nil
}

type GetAttendancePeriodListParam struct {
	StartDate string
	EndDate   string
	Page     int64
	PageSize int64
}

type GetAttendancePeriodListResult struct {
	Data       []model.AttendancePeriod
	Pagination model.Pagination
}

func (store *Store) GetAttendancePeriodList(ctx context.Context, param *GetAttendancePeriodListParam) (*GetAttendancePeriodListResult, error) {
	const op errs.Op = "postgres_store/GetAttendancePeriodList"

	query := `
	SELECT * FROM attendance_period
	WHERE 1=1
	`
	queryParams := []interface{}{}
	paramIdx := 1

	// Filtering by search
	if param.StartDate != "" {
		query += fmt.Sprintf(" AND start_date >= $%d", paramIdx)
		queryParams = append(queryParams, param.StartDate)
		paramIdx++
	}
	if param.EndDate != "" {
		query += fmt.Sprintf(" AND end_date <= $%d", paramIdx)
		queryParams = append(queryParams, param.EndDate)
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
		query += fmt.Sprintf(" ORDER BY start_date DESC LIMIT $%d OFFSET $%d", paramIdx, paramIdx+1)
		queryParams = append(queryParams, param.PageSize, offset)
	} else {
		query += " ORDER BY start_date DESC"
	}

	// Fetch result
	items := []model.AttendancePeriod{}
	if err := store.Db.SelectContext(ctx, &items, query, queryParams...); err != nil {
		store.logger.WithFields(logrus.Fields{
			"op":    op,
			"scope": "Select",
			"err":   err.Error(),
		}).Error()
		return nil, err
	}

	hasNext := param.Page < totalPages

	result := &GetAttendancePeriodListResult{
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
