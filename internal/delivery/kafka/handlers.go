package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/k0kubun/pp"
	pb "product_service/genproto/product"
	"product_service/internal/entity"
	"product_service/internal/infrastructure/kafka"
	"product_service/internal/pkg/config"
	"product_service/internal/usecase"
	"product_service/internal/usecase/event"

	"go.uber.org/zap"
)

type ProductConsumerHandler struct {
	config         *config.Config
	brokerConsumer event.BrokerConsumer
	logger         *zap.Logger
	productUsecase usecase.Product
}

func NewProductConsumerHandler(conf *config.Config,
	brokerConsumer event.BrokerConsumer,
	logger *zap.Logger,
	userUseCase usecase.Product) *ProductConsumerHandler {
	return &ProductConsumerHandler{
		config:         conf,
		brokerConsumer: brokerConsumer,
		logger:         logger,
		productUsecase: userUseCase,
	}
}

func (u *ProductConsumerHandler) HandlerEvents() error {
	pp.Println("1111111111111111111")
	consumerConfig := kafka.NewConsumerConfig(
		u.config.Kafka.Address,
		u.config.Kafka.Topic.ProductService,
		"1",
		func(ctx context.Context, key, value []byte) error {
			var product pb.Product

			if err := json.Unmarshal(value, &product); err != nil {
				return err
			}
			req := entity.Product{
				ID:          product.ID,
				Name:        product.Name,
				Description: product.Description,
				CategoryID:  product.CategoryID,
				UnitPrice:   float64(product.UnitPrice),
				CreatedAt:   product.CreatedAt,
				UpdatedAt:   product.UpdatedAt,
			}
			_, err := u.productUsecase.AddProduct(ctx, &req)
			if err != nil {
				fmt.Println(err, "Add=========================")
			}

			return nil
		},
	)

	u.brokerConsumer.RegisterConsumer(consumerConfig)
	u.brokerConsumer.Run()

	return nil
}
