package services

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	pb "product_service/genproto/product"
	"product_service/internal/delivery/grpc"
	"product_service/internal/entity"
	"product_service/internal/usecase"
)

type productRPC struct {
	logger         *zap.Logger
	productUseCase usecase.Product
	pb.UnimplementedProductServiceServer
	//brokerConsumer event.BrokerConsumer
}

func NewRPC(logger *zap.Logger, productUseCase usecase.Product) *productRPC {
	//brokerProducer event.BrokerProducer) pb.ProductServiceServer {
	return &productRPC{
		logger:         logger,
		productUseCase: productUseCase,
		//brokerConsumer: brokerProducer,
	}
}

func (s productRPC) AddProduct(ctx context.Context, in *pb.Product) (*pb.Product, error) {
	addedProduct, err := s.productUseCase.AddProduct(ctx, &entity.Product{
		ID:          in.ID,
		Name:        in.Name,
		Description: in.Description,
		CategoryID:  in.CategoryID,
		UnitPrice:   float64(in.UnitPrice),
		CreatedAt:   in.CreatedAt,
		UpdatedAt:   in.UpdatedAt,
	})
	if err != nil {
		s.logger.Error("articleCategoriesUsecase.Get", zap.Error(err))
		return &pb.Product{}, grpc.Error(ctx, err)
	}

	return &pb.Product{
		ID:          addedProduct.ID,
		Name:        addedProduct.Name,
		Description: addedProduct.Description,
		CategoryID:  addedProduct.CategoryID,
		UnitPrice:   float32(addedProduct.UnitPrice),
		CreatedAt:   addedProduct.CreatedAt,
		UpdatedAt:   addedProduct.UpdatedAt,
	}, nil
}

func (s productRPC) GetProduct(ctx context.Context, in *pb.IdRequest) (*pb.Product, error) {
	gotProduct, err := s.productUseCase.GetProduct(ctx, in.ID)
	if err != nil {
		s.logger.Error("articleCategoriesUsecase.Get", zap.Error(err))
		return &pb.Product{}, grpc.Error(ctx, err)
	}

	return &pb.Product{
		ID:          gotProduct.ID,
		Name:        gotProduct.Name,
		Description: gotProduct.Description,
		CategoryID:  gotProduct.CategoryID,
		UnitPrice:   float32(gotProduct.UnitPrice),
		CreatedAt:   gotProduct.CreatedAt,
		UpdatedAt:   gotProduct.UpdatedAt,
	}, nil
}

func (s productRPC) UpdateProduct(ctx context.Context, in *pb.Product) (*pb.Product, error) {
	addedProduct, err := s.productUseCase.UpdateProduct(ctx, &entity.Product{
		ID:          in.ID,
		Name:        in.Name,
		Description: in.Description,
		CategoryID:  in.CategoryID,
		UnitPrice:   float64(in.UnitPrice),
		CreatedAt:   in.CreatedAt,
		UpdatedAt:   in.UpdatedAt,
	})
	if err != nil {
		s.logger.Error("articleCategoriesUsecase.Get", zap.Error(err))
		return &pb.Product{}, grpc.Error(ctx, err)
	}

	return &pb.Product{
		ID:          addedProduct.ID,
		Name:        addedProduct.Name,
		Description: addedProduct.Description,
		CategoryID:  addedProduct.CategoryID,
		UnitPrice:   float32(addedProduct.UnitPrice),
		CreatedAt:   addedProduct.CreatedAt,
		UpdatedAt:   addedProduct.UpdatedAt,
	}, nil
}

func (s productRPC) DeleteProduct(ctx context.Context, in *pb.IdRequest) (*pb.Product, error) {
	gotProduct, err := s.productUseCase.DeleteProduct(ctx, in.ID)
	if err != nil {
		s.logger.Error("productCategoriesUsecase.Delete", zap.Error(err))
		return &pb.Product{}, grpc.Error(ctx, err)
	}

	return &pb.Product{
		ID:          gotProduct.ID,
		Name:        gotProduct.Name,
		Description: gotProduct.Description,
		CategoryID:  gotProduct.CategoryID,
		UnitPrice:   float32(gotProduct.UnitPrice),
		CreatedAt:   gotProduct.CreatedAt,
		UpdatedAt:   gotProduct.UpdatedAt,
	}, nil
}
func (s productRPC) GetAllProducts(_ *emptypb.Empty, stream pb.ProductService_GetAllProductsServer) error {
	products, err := s.productUseCase.GetAllProducts(stream.Context())
	if err != nil {
		s.logger.Error("productUseCase.GetAllProducts", zap.Error(err))
		return grpc.Error(stream.Context(), err)
	}

	for _, product := range products {
		if err := stream.Send(&pb.Product{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			CategoryID:  product.CategoryID,
			UnitPrice:   float32(product.UnitPrice),
			CreatedAt:   product.CreatedAt,
			UpdatedAt:   product.UpdatedAt,
		}); err != nil {
			s.logger.Error("stream.Send", zap.Error(err))
			return err
		}
	}

	return nil
}
