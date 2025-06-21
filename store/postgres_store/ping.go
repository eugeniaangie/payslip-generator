package postgres_store

import (
	"context"

	"payslip-generator/util/errs"

	"github.com/sirupsen/logrus"
)

type PingParam struct {
}

type PingResult struct {
	Message string `json:"message"`
}

func (store *Store) Ping(ctx context.Context, param *PingParam) (*PingResult, error) {
	const op errs.Op = "postgres_store/Ping"

	store.logger.WithFields(logrus.Fields{
		"op": op,
	}).Debug("Ping!")

	// init return value
	result := &PingResult{}

	// tidy up result
	result.Message = "pong"

	return result, nil
}
