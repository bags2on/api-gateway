package main

import (
	"context"
	"fmt"
	"time"

	"github.com/bags2on/api-gateway/graph/model"
)

type queryResolver struct {
	server *Server
}

func (r *queryResolver) Products(c context.Context) ([]*model.Product, error) {

	ctx, cancel := context.WithTimeout(c, 3 * time.Second)
	defer cancel()

	productsList, err := r.server.productsClient.GetProducts(ctx)
	if err != nil {
		fmt.Println("cannot get products list")
	}

	var products []*model.Product

	for _, value := range productsList {
		product := &model.Product{
			ID:     value.ID,
			Title:  value.Title,
			Price:  value.Price,
			Images: value.Images,
		}

		products = append(products, product)
	}

	return products, nil
}
