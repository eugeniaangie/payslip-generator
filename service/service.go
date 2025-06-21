package service

import (
	"payslip-generator/store/postgres_store"

	"github.com/sirupsen/logrus"
)

type store struct {
	postgres *postgres_store.Store
}

func newStore(postgresStore *postgres_store.Store) *store {
	return &store{
		postgres: postgresStore,
	}
}

type Service struct {
	logger *logrus.Logger

	store *store
}

func NewService(
	logger *logrus.Logger,
	postgresStore *postgres_store.Store,
) *Service {
	store := newStore(postgresStore)

	return &Service{
		logger: logger,
		store:  store,
	}
}
