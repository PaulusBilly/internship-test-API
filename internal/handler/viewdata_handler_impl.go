package handler

import (
	"internship-test-api/internal/service"
)

func NewViewDataHandler(svc *service.ViewDataService) ViewDataHandler {
	return svc
}
