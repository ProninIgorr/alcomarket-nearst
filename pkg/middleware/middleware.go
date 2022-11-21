package middleware

import (
	"github.com/d-kolpakov/logger"
)

type Middleware struct {
	logger *logger.Logger
}

func New(l *logger.Logger) *Middleware {
	m := &Middleware{}
	m.logger = l
	return m
}
