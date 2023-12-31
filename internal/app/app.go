package app

import (
	"context"
	"os"

	"github.com/elijahelrod/vespene/config"
	"github.com/elijahelrod/vespene/internal/algo/signal"
	"github.com/elijahelrod/vespene/internal/algo/strategy"
	"github.com/elijahelrod/vespene/internal/websocket"
	"github.com/elijahelrod/vespene/pkg/exchange/coinbase"
	"github.com/elijahelrod/vespene/pkg/logger/zap"
)

func Run(ctx context.Context, cfg *config.Config) {

	// Setup Logger Provider (Zap)
	loggerProvider := zap.NewLogger(cfg.Logger.Level, cfg.Logger.DisableCaller, cfg.Logger.DisableStacktrace)
	loggerProvider.InitLogger()

	// Setup DB Client/Provider (MySQL)
	//dbProvider, err := mysql.NewClient(cfg.Database.Host, cfg.Database.User, cfg.Database.Password, cfg.Database.Base, loggerProvider)
	//if err != nil {
	//	loggerProvider.Fatal(err)
	//}
	//defer func() {
	//	err = dbProvider.Close()
	//  if err != nil {
	// loggerProvider.Fatal(err)
	// os.Exit(1)
	// }
	//}()
	//loggerProvider.Info("Initialized DB Client")

	exchangeProvider, err := coinbase.NewClient(cfg.Exchange.Wss) // Websocket Client for Real-Time Market Updates
	if err != nil {
		loggerProvider.Fatal(err)
		return
	}
	defer func() {
		err = exchangeProvider.CloseConnection()
		if err != nil {
			loggerProvider.Fatal(err)
			os.Exit(1)
		}
	}()

	loggerProvider.Info("Initialized Exchange Client")

	socketProvider, err := websocket.NewClient(exchangeProvider, loggerProvider, cfg.Exchange)
	if err != nil {
		loggerProvider.Fatal(err)
	}
	loggerProvider.Info("Initialized Socket Connection To Exchange: ", cfg.Exchange.Wss)

	testStrategy := strategy.NewStrategy("Test", loggerProvider, signal.NewDonchian(5, 4))

	go func() {
		loggerProvider.Info("Socket starting...")
		if err = socketProvider.Run(ctx, *testStrategy); err != nil {
			loggerProvider.Fatal(err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()

	loggerProvider.Info("Socket stopping...")
}
