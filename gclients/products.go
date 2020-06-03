package gclients

import (
	"context"

	"github.com/bags2on/api-gateway/graph/model"
	proto "github.com/bags2on/api-gateway/protobuf"
	"google.golang.org/grpc"
)

type ProductsClient struct {
	connection *grpc.ClientConn
	service    proto.ProductServiceClient
}

func GetProductsClient(url string) (*ProductsClient, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := proto.NewProductServiceClient(conn)

	return &ProductsClient{conn, client}, nil
}

func (p *ProductsClient) Close() {
	p.connection.Close()
}

func (p *ProductsClient) GetProducts(ctx context.Context) ([]model.Product, error) {
	r, err := p.service.GetProducts(ctx, &proto.ProductRequest{})
	if err != nil {
		return nil, err
	}

	var products []model.Product

	for _, value := range r.Products {
		products = append(products, model.Product{
			ID:      value.Id,
			Title:   value.Title,
			Price:   int(value.Price),
			Preview: value.Preview,
		})
	}

	return products, nil
}

func (p *ProductsClient) GetProductByID(ctx context.Context, id string) (*model.Product, error) {
	r, err := p.service.GetProductByID(ctx, &proto.ProductByIdRequest{
		Id: id,
	})

	if err != nil {
		return nil, err
	}

	return &model.Product{
		ID:          r.Product.Id,
		Title:       r.Product.Title,
		Price:       int(r.Product.Price),
		Images:      r.Product.Images,
		Preview:     r.Product.Preview,
		Tags:        r.Product.Tags,
		Description: &r.Product.Description,
	}, nil
}
