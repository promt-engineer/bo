package handlers

import (
	"backoffice/internal/constants"
	"backoffice/internal/services"
	"backoffice/internal/transport/queue"
	"context"
	"go.uber.org/zap"
)

type currencyHandler struct {
	configSenderService *services.ConfigSenderService
}

func NewCurrencyHandler(configSenderService *services.ConfigSenderService) *currencyHandler {
	return &currencyHandler{configSenderService: configSenderService}
}

func (c currencyHandler) Register(r *queue.Router) {
	r.Accept(constants.MsgCurrencyRequestType, func(body []byte) {
		zap.S().Info("get request from lord for currency")
		if err := c.configSenderService.SendCurrencyToLord(context.Background()); err != nil {
			zap.S().Error(err)
		}
	})
}
