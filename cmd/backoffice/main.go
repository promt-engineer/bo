//go:build !codeanalysis
// +build !codeanalysis

package main

import (
	"context"
	"log"
	"sync"
	"time"

	"backoffice/internal/constants"
	"backoffice/internal/container"
	"backoffice/internal/transport/http"
	"backoffice/internal/transport/queue"
	"backoffice/internal/transport/rpc"
	"backoffice/pkg/totp"
	"backoffice/pkg/validator"
	"backoffice/utils"
	"bitbucket.org/play-workspace/gocommon/tracer"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"
)

func main() {
	now := time.Now()
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	app := container.Build(ctx, wg)
	_ = app.Get(constants.LoggerName).(*zap.Logger)
	_ = app.Get(constants.TOTPName).(*totp.TOTP)

	zap.S().Info("Starting application...")
	zap.S().Infof("Up and running (%s)", time.Since(now))

	tr := app.Get(constants.TracerName).(*tracer.JaegerTracer)
	queueManager := app.Get(constants.QueueName).(*queue.Queue)
	queueManager.AddHandlers(app.Get(constants.CurrencyQueueHandlerName).(queue.Handler))
	queueManager.WaitTillStart()

	server := app.Get(constants.HTTPServerName).(*http.Server)
	binding.Validator = app.Get(constants.ValidatorName).(*validator.Validator)

	go server.Run()

	rpcHandler := app.Get(constants.RPCName).(*rpc.Handler)

	go rpc.StartUnsecureRPCServer(rpcHandler, tr)

	zap.S().Infof("Up and running (%s)", time.Since(now))
	zap.S().Infof("Got %s signal. Shutting down...", <-utils.WaitTermSignal())

	if err := app.Delete(); err != nil {
		log.Println(err)
	}

	cancel()
	wg.Wait()

	zap.S().Info("Service stopped.")
}
