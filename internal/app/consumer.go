//package app
//
//import (
//	"go.uber.org/zap"
//	"product_service/internal/infrastructure/kafka"
//	"product_service/internal/pkg/config"
//	logpkg "product_service/internal/pkg/logger"
//	"product_service/internal/pkg/postgres"
//	"product_service/internal/usecase/event"
//)
//
//type ProductConsumer struct {
//	Config         *config.Config
//	Logger         *zap.Logger
//	DB             *postgres.PostgresDB
//	BrokerConsumer event.BrokerConsumer
//}
//
//func NewProductConsumer(conf *config.Config) (*ProductConsumer, error) {
//	logger, err := logpkg.New(conf.LogLevel, conf.Environment, conf.APP+"_cli"+".lo")
//	if err != nil {
//		return nil, err
//	}
//
//	consumer := kafka.NewConsumer(logger)
//
//	db, err := postgres.New(conf)
//	if err != nil {
//		return nil, err
//	}
//
//}

package app

import (
	"go.uber.org/zap"
	handlers "product_service/internal/delivery/kafka"
	"product_service/internal/infrastructure/kafka"
	"product_service/internal/infrastructure/repository/postgresql"
	"product_service/internal/pkg/config"
	logpkg "product_service/internal/pkg/logger"
	"product_service/internal/pkg/postgres"
	"product_service/internal/usecase"
	"product_service/internal/usecase/event"
	"time"
)

type UserConsumer struct {
	Config         *config.Config
	Logger         *zap.Logger
	DB             *postgres.PostgresDB
	BrokerConsumer event.BrokerConsumer
}

func NewProductConsumer(conf *config.Config) (*UserConsumer, error) {
	logger, err := logpkg.New(conf.LogLevel, conf.Environment, conf.APP+"_cli"+".lo")
	if err != nil {
		return nil, err
	}

	consumer := kafka.NewConsumer(logger)

	db, err := postgres.New(conf)
	if err != nil {
		return nil, err
	}

	return &UserConsumer{Config: conf, Logger: logger, DB: db, BrokerConsumer: consumer}, nil
}

func (u *UserConsumer) Run() error {

	// repo init
	userRepo := postgresql.NewProductRepo(u.DB)

	// usecase init
	ctx, err := time.ParseDuration(u.Config.Context.Timeout)
	if err != nil {
		return err
	}
	productUseCase := usecase.NewProductService(ctx, userRepo)

	// event handler
	eventHandler := handlers.NewProductConsumerHandler(u.Config, u.BrokerConsumer, u.Logger, productUseCase)

	return eventHandler.HandlerEvents()
}

func (u *UserConsumer) Close() {
	u.BrokerConsumer.Close()

	u.Logger.Sync()
}
