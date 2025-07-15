package framework

import (
    "go.uber.org/zap"
)

// NewLogger returns a production-ready SugaredLogger
func NewLogger() (*zap.SugaredLogger, error) {
    lg, err := zap.NewProduction()
    if err != nil {
        return nil, err
    }
    return lg.Sugar(), nil
}