package controller

import (
	"payslip-generator/service"
	"time"
)

// controller setup

type setup struct {
	timeout time.Duration
}

func newSetup(timeout time.Duration) *setup {
	setup := &setup{
		timeout: timeout,
	}

	return setup
}

type Controller struct {
	setup *setup

	service *service.Service
}

func NewController(timeout time.Duration, service *service.Service) *Controller {
	setup := newSetup(timeout)

	return &Controller{
		setup: setup,

		service: service,
	}
}
