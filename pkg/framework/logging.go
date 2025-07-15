package framework

import (
    "go.uber.org/zap"
)

// NewLogger returns a SugaredLogger with production config
func NewLogger() (*zap.SugaredLogger, error) {
    logger, err := zap.NewProduction()
    if err != nil {
        return nil, err
    }
    return logger.Sugar(), nil
}
